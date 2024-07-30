package handlers

import (
	"WebScraper/models"
	"WebScraper/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetScrapedDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, err := os.Open("scraped_data.json")
	if err != nil {
		http.Error(w, "Failed to open scraped data file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read scraped data file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Scrape and extract links from the page
func ScrapeAndExtractLinks(pageURL string) ([]string, error) {
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

	finalText := strings.TrimSpace(textContent.String())
	finalText = utils.RemoveBlankLines(finalText)
	finalText = utils.RemoveExtraSpaces(finalText)

	pageTitle := doc.Find("title").Text()

	page := models.PageData{
		Title:   pageTitle,
		URL:     pageURL,
		Content: finalText,
	}

	appendToFile(page)

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
func appendToFile(page models.PageData) {
	fileName := "scraped_data.json"
	var pages []models.PageData

	if _, err := os.Stat(fileName); err == nil {
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

	pages = append(pages, page)

	jsonData, err := json.MarshalIndent(pages, "", "    ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

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
