package handlers

import (
	"WebScraper/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	var requestData models.URLDatastruct // getting the model from the models package

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Closing the request body, it's a must
	defer r.Body.Close()

	jobs := make(chan string, len(requestData.URLs))
	results := make(chan struct{}, len(requestData.URLs))

	// Start a single worker to process URLs sequentially
	wg.Add(1)
	go worker(1, jobs, results)

	// Send URLs to the jobs channel sequentially
	go func() {
		defer close(jobs)
		for _, url := range requestData.URLs {
			jobs <- url
		}
	}()

	// Waiting for all work to be done
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
			Crawl(url) // Process the URL fully before moving to the next
			results <- struct{}{}
		}(url)
	}
}

// Crawl function to check all the "visited" pages
func Crawl(baseURL string) {
	toVisit := []string{baseURL}

	for len(toVisit) > 0 {
		url := toVisit[0]
		toVisit = toVisit[1:]

		visitLock.Lock()
		if visited[url] {
			visitLock.Unlock()
			continue
		}
		visited[url] = true
		visitLock.Unlock()

		fmt.Println("Fetching:", url)

		links, err := ScrapeAndExtractLinks(url)
		if err != nil {
			// Log the error but continue processing this URL fully
			log.Printf("Error scraping %s: %v\n", url, err)
			continue
		}

		for _, link := range links {
			visitLock.Lock()
			if !visited[link] {
				toVisit = append(toVisit, link)
			}
			visitLock.Unlock()
		}
	}
}
