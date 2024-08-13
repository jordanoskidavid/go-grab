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
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout for the context
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var pageTitle, textContent string

	err := chromedp.Run(ctx,
		network.SetBlockedURLS([]string{"*.jpg", "*.png", "*.gif", "*.css", "*.svg", "*.js"}),
		chromedp.Navigate(pageURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Title(&pageTitle),
		chromedp.Evaluate(`document.body.innerText`, &textContent),
	)
	if err != nil {
		return nil, fmt.Errorf("error rendering dynamic content from %s: %v", pageURL, err)
	}

	finalText := strings.TrimSpace(textContent)
	finalText = utils.RemoveBlankLines(finalText)
	finalText = utils.RemoveExtraSpaces(finalText)

	pageData := models.PageData{
		Title:   pageTitle,
		URL:     pageURL,
		Content: finalText,
	}

	if err := savePageToFile(pageData); err != nil {
		return nil, err
	}

	var links []string
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a[href^="http"]')).map(a => a.href)`, &links),
	)
	if err != nil {
		return nil, fmt.Errorf("error extracting links from %s: %v", pageURL, err)
	}

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

func savePageToFile(page models.PageData) error {
	baseURL, err := getBaseURL(page.URL)
	if err != nil {
		return err
	}
	fileName := sanitizeFileName(baseURL) + ".json"
	filePath := filepath.Join(".", fileName)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var pages []models.PageData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pages); err != nil && err != io.EOF {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	pages = append(pages, page)

	file.Truncate(0)
	file.Seek(0, 0)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(pages); err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}
	fmt.Printf("Page data saved to: %s\n", absPath)
	return nil
}

func getBaseURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}
	return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host), nil
}

func sanitizeFileName(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("error parsing URL: %v", err)
		return "invalid_url"
	}
	fileName := parsedURL.Hostname()
	if len(fileName) == 0 {
		fileName = "default"
	}
	return fileName
}
