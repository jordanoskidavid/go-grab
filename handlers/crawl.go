package handlers

import (
	"WebScraper/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	visited   = make(map[string]bool)
	visitLock sync.Mutex
	wg        sync.WaitGroup
)

func StartCrawlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData models.URLDatastruct

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	jobs := make(chan string, len(requestData.URLs))
	results := make(chan struct{}, len(requestData.URLs))

	numWorkers := 5 // Adjust this number based on your needs
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results)
	}

	go func() {
		defer close(jobs)
		for _, url := range requestData.URLs {
			jobs <- url
		}
	}()

	go func() {
		wg.Wait()
		close(results)
		fmt.Println("Crawling completed.")
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Crawling started"))
}

func worker(id int, jobs <-chan string, results chan<- struct{}) {
	defer wg.Done()
	for url := range jobs {
		func(url string) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Worker %d: Recovered from panic while processing %s: %v", id, url, r)
				}
			}()
			Crawl(url)
			results <- struct{}{}
		}(url)
	}
}

func normalizeURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	path := strings.TrimRight(u.Path, "/")
	u.Path = path
	return u.String()
}
func Crawl(baseURL string) {
	toVisit := []string{baseURL}

	for len(toVisit) > 0 {
		url := toVisit[0]
		toVisit = toVisit[1:]

		// Normalize URL before checking and updating the visited map
		normalizedURL := normalizeURL(url)

		visitLock.Lock()
		if visited[normalizedURL] {
			visitLock.Unlock()
			continue
		}
		visited[normalizedURL] = true
		visitLock.Unlock()

		fmt.Println("Fetching:", url)

		links, err := ScrapeAndExtractLinks(url)
		if err != nil {
			log.Printf("Error scraping %s: %v\n", url, err)
			continue
		}

		for _, link := range links {
			normalizedLink := normalizeURL(link)
			visitLock.Lock()
			if !visited[normalizedLink] {
				toVisit = append(toVisit, link)
			}
			visitLock.Unlock()
		}
	}
}
