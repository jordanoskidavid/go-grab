package handlers

import (
	"WebScraper/models"
	"WebScraper/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
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
func ScrapeAndExtractLinks(pageURL string) ([]string, error) {
	// Create a new Chrome context with a timeout
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout for the context
	ctx, cancel = context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var htmlContent string
	var pageTitle string

	// Set blocked URLs and run the chromedp tasks to navigate and fetch the rendered HTML and title
	err := chromedp.Run(ctx,
		network.SetBlockedURLS([]string{"*.jpg", "*.png", "*.gif", "*.css", "*.svg"}), // Block specific resources
		chromedp.Navigate(pageURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent),
		chromedp.Title(&pageTitle),
	)
	if err != nil {
		return nil, fmt.Errorf("error rendering dynamic content from %s: %v", pageURL, err)
	}

	// Extract human-readable text using chromedp
	var textContent string
	err = chromedp.Run(ctx,
		chromedp.Text("body", &textContent, chromedp.NodeVisible),
	)
	if err != nil {
		return nil, fmt.Errorf("error extracting text from %s: %v", pageURL, err)
	}

	// Clean and format the text content
	finalText := strings.TrimSpace(textContent)
	finalText = utils.RemoveBlankLines(finalText)
	finalText = utils.RemoveExtraSpaces(finalText)

	// Save the page content to a file
	pageData := models.PageData{
		Title:   pageTitle,
		URL:     pageURL,
		Content: finalText,
	}

	err = savePageToFile(pageData)
	if err != nil {
		return nil, err
	}

	// Extract and filter links
	var links []string
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href).filter(href => href.startsWith('http'))`, &links),
	)
	if err != nil {
		return nil, fmt.Errorf("error extracting links from %s: %v", pageURL, err)
	}

	// Filter out external links
	base, err := url.Parse(pageURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL %s: %v", pageURL, err)
	}
	baseHost := base.Hostname()

	var internalLinks []string
	for _, link := range links {
		parsedLink, err := url.Parse(link)
		if err != nil {
			continue
		}
		if parsedLink.Hostname() == baseHost {
			internalLinks = append(internalLinks, link)
		}
	}

	return internalLinks, nil
}

// savePageToFile saves the page content to a file named after the base URL.
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

	// Append new page data
	pages = append(pages, page)

	// Write updated data back to the file
	jsonData, err := json.MarshalIndent(pages, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
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
