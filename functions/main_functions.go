package functions

import (
	"GoGrab/models"
	"GoGrab/utils"
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var (
	visited   = make(map[string]bool) // keeps track of visited URLs
	visitLock sync.Mutex              // here using mutex i.e ensuring that only one goroutine at a time to access a shared resource at him
)

/*
	 	Crawl initiates a web crawling process starting from the base URL. It visits the URL, extracts links from it, and adds
		them to the visit queue if the haven't been visited yet. The function also uses locking to ensure thread-safe operations
		when managing list of visited URLs.
*/
func Crawl(baseURL string) {
	toVisit := []string{baseURL}     // List of URLs to visit
	visitLock := sync.Mutex{}        // Mutex to ensure safe access to the visited map
	visited := make(map[string]bool) // Tracks visited URLs
	var wg sync.WaitGroup            // WaitGroup to manage concurent goroutines

	for len(toVisit) > 0 {
		url := toVisit[0]     // Get the next URL to visit
		toVisit = toVisit[1:] // Remove the URL from the visit list

		// Adding a delay to avoid overloading the target site
		time.Sleep(1 * time.Second)

		// Normalize the URL to ensure consistent comparisons
		normalizedURL := utils.NormalizeURL(url)

		// Lock the visitLock to check and update the visited map safely
		visitLock.Lock()
		if visited[normalizedURL] {
			visitLock.Unlock()
			continue // Skip the URL if it has already been visited
		}
		visited[normalizedURL] = true // Mark the URL as visited
		visitLock.Unlock()

		// Log the fetching process
		fmt.Println("Fetching:", url)

		// Scrape the URL and extract links from the page
		links, err := ScrapeAndExtractLinks(url)
		if err != nil {
			// Log the error if scraping fails
			log.Printf("Error scraping %s: %v\n", url, err)
			continue
		}

		// Process each extracted link
		for _, link := range links {
			normalizedLink := utils.NormalizeURL(link)
			visitLock.Lock()
			if !visited[normalizedLink] {
				toVisit = append(toVisit, link) // Add new links to the visit list
			}
			visitLock.Unlock()
		}
	}
	// Wait for all goroutines to finish before exiting
	wg.Wait()
}

/*
ScrapeAndExtractLinks scrapes a given page URL, extracts its content and internal links.
It uses Chrome DevTools Protocol (CDP) to navigate the page, block unnecessary assets, and extract both text and links.
*/
func ScrapeAndExtractLinks(pageURL string) ([]string, error) {
	// Create a new browser context for the scraping task
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout for the scraping operation
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var pageTitle, textContent string

	// Run Chrome DevTools Protocol (CDP) tasks to block unnecessary assets, navigate the page, and extract the page title and body text
	err := chromedp.Run(ctx,
		network.SetBlockedURLS([]string{"*.jpg", "*.png", "*.gif", "*.css", "*.svg", "*.js"}),
		chromedp.Navigate(pageURL),                                 //navigate to the page
		chromedp.WaitVisible("body", chromedp.ByQuery),             // Wait until the body is visible
		chromedp.Title(&pageTitle),                                 // Extract the page title
		chromedp.Evaluate(`document.body.innerText`, &textContent), // Extract the body text content
	)
	if err != nil {
		// Return an error if any of the scraping steps fail
		return nil, fmt.Errorf("error rendering dynamic content from %s: %v", pageURL, err)
	}

	// Format the extracted text content by removing blank lines and extra spaces
	finalText := strings.TrimSpace(textContent)
	finalText = utils.RemoveBlankLines(finalText)
	finalText = utils.RemoveExtraSpaces(finalText)

	// Create a PageData model and save the extracted data to a file
	pageData := models.PageData{
		Title:   pageTitle,
		URL:     pageURL,
		Content: finalText,
	}
	// Save the scraped page content to a file using the utility function
	if err := utils.SavePageToFile(pageData); err != nil {
		return nil, err
	}

	// Extract all links (<a> tags) from the page using CDP
	var links []string
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a[href]')).map(a => a.href)`, &links),
	)
	if err != nil {
		// Return an error if link extraction fails
		return nil, fmt.Errorf("error extracting links from %s: %v", pageURL, err)
	}
	// Parse the base URL to extract the hostname
	base, err := url.Parse(pageURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL %s: %v", pageURL, err)
	}
	baseHost := base.Hostname()

	// Filter and return only internal links (i.e., links within the same domain)
	var internalLinks []string
	for _, link := range links {
		parsedLink, err := url.Parse(link)
		if err != nil {
			continue //skip invalid links
		}
		if parsedLink.Hostname() == baseHost {
			internalLinks = append(internalLinks, link) // Add internal links to the result
		}
	}

	return internalLinks, nil // Return the list of internal linkss
}
