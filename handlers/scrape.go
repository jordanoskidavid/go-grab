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

	"github.com/PuerkitoBio/goquery"
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
	// Create a context with a timeout to manage the Chrome session
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var htmlContent string

	// Run the chromedp tasks to navigate and fetch the rendered HTML
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return nil, fmt.Errorf("error rendering dynamic content from %s: %v", pageURL, err)
	}

	// Use goquery to parse the fetched HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("error loading HTML for URL %s: %v", pageURL, err)
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

func savePageToFile(page models.PageData) error {
	baseURL, err := getBaseURL(page.URL)
	if err != nil {
		return err
	}
	fileName := sanitizeFileName(baseURL) + ".json"
	filePath := filepath.Join(".", fileName)
	var pages []models.PageData
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

	pages = append(pages, page)

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

func getBaseURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}
	baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	return baseURL, nil
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
