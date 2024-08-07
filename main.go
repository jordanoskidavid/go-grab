package main

import (
	"WebScraper/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/crawl", handlers.StartCrawlHandler)
	http.HandleFunc("/scraped-data", handlers.GetScrapedDataHandler)
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
