package handlers

import (
	"WebScraper/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var visited = make(map[string]bool)

func StartCrawlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData models.URLDatastruct // Use the model from the models package

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Close the request body
	defer r.Body.Close()

	// Process each URL
	for _, url := range requestData.URLs {
		go Crawl(url) // Run crawling in a goroutine
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Crawling started"))
}

// Crawl function to check all the "visited" pages
func Crawl(baseURL string) {
	toVisit := []string{baseURL}

	for len(toVisit) > 0 {
		url := toVisit[0]
		toVisit = toVisit[1:]
		if visited[url] {
			continue
		}

		fmt.Println("Fetching:", url)
		visited[url] = true

		links, err := ScrapeAndExtractLinks(url)
		if err != nil {
			log.Printf("Error scraping %s: %v\n", url, err)
			continue
		}

		for _, link := range links {
			if !visited[link] {
				toVisit = append(toVisit, link)
			}
		}
	}
}
