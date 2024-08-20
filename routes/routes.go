package routes

import (
	"WebScraper/handlers"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/api/crawl", handlers.StartCrawlHandler)
	http.HandleFunc("/api/get-data", handlers.GetScrapedDataHandler)
	http.HandleFunc("/api/delete-data", handlers.DeleteScrapedData)
	http.HandleFunc("/api/generate-api-key", handlers.GenerateAPIKeyHandler)
}
