package handlers

import (
	"WebScraper/functions"
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

	urls, err := functions.Read_json_urls("urls.json")
	if err != nil {
		http.Error(w, "Failed to read URLs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, url := range urls {
		go Crawl(url) // Start crawling URLs in the background
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
