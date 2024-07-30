package main

import (
	"WebScraper/functions"
	"WebScraper/handlers"
	"fmt"
)

func main() {
	urls, err := functions.Read_json_urls("urls.json")
	if err != nil {
		fmt.Println("The reading can not be done:", err)
		return
	}
	for _, url := range urls {
		handlers.Crawl(url) // Crawling the URL
	}
	fmt.Println("Crawling completed.")
}
