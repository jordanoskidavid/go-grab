package handlers

import (
	"WebScraper/functions"
	"WebScraper/models"
	"encoding/json"
	"net/http"
)

// StartCrawlHandler godoc
// @Summary Starts a web crawl process
// @Description Initiates a web scraping process by accepting a list of URLs.
// @Tags Crawling
// @Accept json
// @Produce json
// @Param request body models.URLDatastruct true "List of URLs to crawl"
// @Success 200 {string} string "Crawling completed"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 405 {string} string "Invalid request method"
// @Router /api/crawl [post]

func StartCrawlHandler(w http.ResponseWriter, r *http.Request) {
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
		functions.Crawl(url)
	}

	//w.WriteHeader(http.StatusOK)
	w.Write([]byte("Crawling completed"))
}
