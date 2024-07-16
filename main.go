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

var visited = make(map[string]bool)

func main() {
	startURL := "https://scrapeme.live/shop/"
	crawl(startURL)
}

func crawl(baseURL string) {
	toVisit := []string{baseURL}

	for len(toVisit) > 0 {
		url := toVisit[0]
		toVisit = toVisit[1:]

		if visited[url] {
			continue
		}

		fmt.Println("Fetching:", url)
		visited[url] = true

		links, err := scrapeAndExtractLinks(url)
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

	doc.Find("style, script, .jquery-script").Remove()

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

	var links []string
	base, err := url.Parse(pageURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL %s: %v", pageURL, err)
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			absLink := base.ResolveReference(&url.URL{Path: link}).String()
			links = append(links, absLink)
		}
	})

	return links, nil
}

func sanitizeFilename(url string) string {
	return strings.ReplaceAll(url, "/", "_")
}
