package main

import (
	"WebScraper/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/crawl", handlers.StartCrawlHandler)
	http.HandleFunc("/get-data", handlers.GetScrapedDataHandler)
	http.HandleFunc("/delete-data", handlers.DeleteScrapedData)
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
