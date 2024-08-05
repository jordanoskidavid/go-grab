package handlers

import (
	"WebScraper/models"
	"WebScraper/utils"
	"crypto/sha256"
	"encoding/hex"
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

// GetScrapedDataHandler serves the scraped data from a file.
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

// ScrapeAndExtractLinks scrapes content from the given URL and extracts links.
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

	// Remove unwanted elements like style, script, etc.
	doc.Find("style, script, .jquery-script").Remove()

	var textContent strings.Builder
	lastWasSpace := true

	// Filter for main content areas only
	filterMainContent := func(i int, s *goquery.Selection) bool {
		tagName := strings.ToLower(s.Get(0).Data)
		return tagName == "p" || tagName == "div" || tagName == "span" || tagName == "a"
	}

	doc.Find("body *").FilterFunction(filterMainContent).Each(func(i int, s *goquery.Selection) {
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

	// Add more logging to inspect content and deduplication process
	fmt.Printf("Scraped content for URL: %s\n", page.URL)
	fmt.Printf("Content hash: %s\n", hashContent(page.Content))

	err = savePageToFile(page)
	if err != nil {
		return nil, err
	}

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

// savePageToFile saves or appends page data to a JSON file named after the base URL.
func savePageToFile(page models.PageData) error {
	baseURL, err := getBaseURL(page.URL)
	if err != nil {
		return err
	}

	fileName := sanitizeFileName(baseURL) + ".json"
	filePath := filepath.Join(".", fileName)

	var pages []models.PageData

	// Check if file exists and read existing data
	if _, err := os.Stat(filePath); err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&pages); err != nil && err != io.EOF {
			return fmt.Errorf("error decoding JSON: %v", err)
		}
	}

	// Deduplication: Check for duplicate content before saving
	contentHash := hashContent(page.Content)
	fmt.Printf("Checking for duplicate content hash: %s\n", contentHash)
	for _, p := range pages {
		if hashContent(p.Content) == contentHash {
			fmt.Printf("Duplicate content detected for URL %s, skipping save.\n", page.URL)
			return nil
		}
	}

	// Append new page data only if not duplicate
	pages = append(pages, page)

	// Write updated data back to the file
	jsonData, err := json.MarshalIndent(pages, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	// Overwrite the file with the new content
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}
	fmt.Printf("Page data saved to: %s\n", absPath)

	return nil
}

// getBaseURL extracts the base URL from a given URL.
func getBaseURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}

	// Construct the base URL
	baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	return baseURL, nil
}

// sanitizeFileName sanitizes a URL to create a valid filename.
func sanitizeFileName(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("error parsing URL: %v", err)
		return "invalid_url"
	}

	// Create a valid file name from the URL
	fileName := parsedURL.Hostname()
	if len(fileName) == 0 {
		fileName = "default"
	}
	return fileName
}

// hashContent generates a hash for the given content to help with deduplication.
func hashContent(content string) string {
	h := sha256.New()
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}
