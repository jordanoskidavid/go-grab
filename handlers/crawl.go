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

	// Process each URL
	for _, url := range requestData.URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			Crawl(url) // Using go routine it runs crawling through the sites
		}(url)
	}

	go func() {
		wg.Wait()
		fmt.Println("Crawling completed.")
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Crawling started"))
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
