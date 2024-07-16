package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Set to keep track of visited URLs
var visited = make(map[string]bool)

func main() {
	startURL := "https://scrapeme.live/shop/"
	crawl(startURL)
}

func crawl(baseURL string) {
	toVisit := []string{baseURL}

	for len(toVisit) > 0 {
		// Get the next URL to visit
		url := toVisit[0]
		toVisit = toVisit[1:]

		// Skip if already visited
		if visited[url] {
			continue
		}

		fmt.Println("Fetching:", url)
		visited[url] = true

		// Scrape and extract links
		links, err := scrapeAndExtractLinks(url)
		if err != nil {
			log.Printf("Error scraping %s: %v\n", url, err)
			continue
		}

		// Add new links to the to-visit list
		for _, link := range links {
			if !visited[link] {
				toVisit = append(toVisit, link)
			}
		}
	}
}

func scrapeAndExtractLinks(pageURL string) ([]string, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL %s: %v", pageURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %d for URL %s", resp.StatusCode, pageURL)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading HTML from URL %s: %v", pageURL, err)
	}

	// Remove <style> and <script> tags
	doc.Find("style, script, .jquery-script").Remove()

	// Extract human-readable text content
	var textContent strings.Builder
	lastWasSpace := true

	filterNonText := func(i int, s *goquery.Selection) bool {
		tagName := strings.ToLower(s.Get(0).Data)
		return tagName == "p" || tagName == "div" || tagName == "span" || tagName == "a"
	}

	doc.Find("body *").FilterFunction(filterNonText).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			if !lastWasSpace && textContent.Len() > 0 {
				textContent.WriteString(" ")
			}
			textContent.WriteString(text)
			lastWasSpace = false
		} else {
			lastWasSpace = true
		}
	})

	// Save extracted text content to a text file
	fileName := fmt.Sprintf("extracted_content_%s.txt", sanitizeFilename(pageURL))
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(textContent.String())
	if err != nil {
		return nil, fmt.Errorf("error writing to file: %v", err)
	}

	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %v", err)
	}
	fmt.Printf("Extracted text content saved to: %s\n", absPath)

	// Extract links
	var links []string
	base, err := url.Parse(pageURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL %s: %v", pageURL, err)
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			// Resolve relative URLs
			absLink := base.ResolveReference(&url.URL{Path: link}).String()
			links = append(links, absLink)
		}
	})

	return links, nil
}

func sanitizeFilename(url string) string {
	return strings.ReplaceAll(url, "/", "_")
}
