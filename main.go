package main

import (
	"encoding/json"
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

// structure for storing the json data
type PageData struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

func main() {
	startURL := "https://scrapeme.live/shop/"
	crawl(startURL) // Crawling the URL
	fmt.Println("Crawling completed.")
}

// Crawl function to check all the "visited" pages
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

// Scrape and extract links from the page
func scrapeAndExtractLinks(pageURL string) ([]string, error) {
	resp, err := http.Get(pageURL) //this function is able to get the url's from the crawl function and check if the fetching url is OK or throws an error and gets the response code back
	if err != nil {
		return nil, fmt.Errorf("error fetching URL %s: %v", pageURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %d for URL %s", resp.StatusCode, pageURL)
	}
	//down here the goquery library gets the html code and check if html is taken and pass it to the go query
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading HTML from URL %s: %v", pageURL, err)
	}
	//here with using doc.find is finding css js and jquery and removes it
	doc.Find("style, script, .jquery-script").Remove()

	var textContent strings.Builder
	lastWasSpace := true

	//down here is filtering the usage of the nontext strings and returns only p div span and a tags

	filterNonText := func(i int, s *goquery.Selection) bool {
		tagName := strings.ToLower(s.Get(0).Data)
		return tagName == "p" || tagName == "div" || tagName == "span" || tagName == "a"
	}
	//here it goes through body page
	doc.Find("body *").FilterFunction(filterNonText).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			text = strings.ReplaceAll(text, "\t", " ")

			if !lastWasSpace && textContent.Len() > 0 {
				textContent.WriteString(" ")
			}
			textContent.WriteString(text)
			lastWasSpace = false
		} else {
			lastWasSpace = true
		}
	})

	// Removing spaces that aint needed
	finalText := strings.TrimSpace(textContent.String())
	finalText = removeBlankLines(finalText)
	finalText = removeExtraSpaces(finalText)

	// Get the page title
	pageTitle := doc.Find("title").Text()

	// Construct the page data
	page := PageData{
		Title:   pageTitle,
		URL:     pageURL,
		Content: finalText,
	}

	// Append the page data to the JSON file
	appendToFile(page)

	// Extract links from the page
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

// Append page data to the JSON file
func appendToFile(page PageData) {
	fileName := "scraped_data.json"
	var pages []PageData

	// Check if the file already exists
	if _, err := os.Stat(fileName); err == nil {
		// File exists, read the current contents
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&pages); err != nil {
			log.Fatalf("Error decoding JSON: %v", err)
		}
	}

	// Append the new page data
	pages = append(pages, page)

	// Marshal the pages slice to JSON
	jsonData, err := json.MarshalIndent(pages, "", "    ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	// Write the JSON data to file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	absPath, err := filepath.Abs(fileName)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}
	fmt.Printf("Page data appended to: %s\n", absPath)
}

// Filter functions
func removeBlankLines(text string) string {
	var result strings.Builder
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			if result.Len() > 0 {
				result.WriteString("\n")
			}
			result.WriteString(trimmed)
		}
	}
	return result.String()
}

func removeExtraSpaces(text string) string {
	words := strings.Fields(text)
	return strings.Join(words, " ")
}
