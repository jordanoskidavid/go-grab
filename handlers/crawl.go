package handlers

import (
	"WebScraper/models"
	"WebScraper/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	visited   = make(map[string]bool)
	visitLock sync.Mutex
)

func StartCrawlHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Crawling started"))

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

	for _, url := range requestData.URLs {
		Crawl(url)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Crawling completed"))
}

func Crawl(baseURL string) {
	toVisit := []string{baseURL}

	for len(toVisit) > 0 {
		url := toVisit[0]
		toVisit = toVisit[1:]

		normalizedURL := utils.NormalizeURL(url)

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
			normalizedLink := utils.NormalizeURL(link)
			visitLock.Lock()
			if !visited[normalizedLink] {
				toVisit = append(toVisit, link)
			}
			visitLock.Unlock()
		}
	}
}
