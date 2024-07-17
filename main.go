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

// making structure that can be used for storing JSON
type PageData struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

func main() {
	startURL := "https://scrapeme.live/shop/"
	crawl(startURL) //crawling the url
}

// setting up the function where checks all the "visited" pages
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

// getting the html code and then transforming it with the goquery library
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

	//Removing the blank spaces
	finalText := strings.TrimSpace(textContent.String())

	//Removing the blank lines
	finalText = removeBlankLines(finalText)

	//Removing spaces that are extrra
	finalText = removeExtraSpaces(finalText)

	//Getting the page title
	pageTitle := doc.Find("title").Text()

	//Constucting the page data where it takes page title, url and filtered text that can be stored
	page := PageData{
		Title:   pageTitle,
		URL:     pageURL,
		Content: finalText,
	}

	// Marshal PageData struct to JSON
	jsonData, err := json.MarshalIndent(page, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Write JSON data to file
	fileName := fmt.Sprintf("scraped_files/scraped_data_%s.json", sanitizeFilename(pageURL))
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return nil, fmt.Errorf("error writing to file: %v", err)
	}

	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %v", err)
	}
	fmt.Printf("Page data saved to: %s\n", absPath)

	//Converting the struct into JSON
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

// These all are filter functions that are used up above :)
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

func sanitizeFilename(url string) string {
	return strings.ReplaceAll(url, "/", "_")
}
