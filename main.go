package main

import (
	"WebScraper/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/crawl", handlers.StartCrawlHandler)
	http.HandleFunc("/scraped-data", handlers.GetScrapedDataHandler) //da se sredi
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
