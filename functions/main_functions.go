package functions

import (
	"WebScraper/models"
	"WebScraper/utils"
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
	visited   = make(map[string]bool)
	visitLock sync.Mutex
)

func Crawl(baseURL string) {
	toVisit := []string{baseURL}
	visitLock := sync.Mutex{}
	visited := make(map[string]bool)
	var wg sync.WaitGroup

	for len(toVisit) > 0 {
		url := toVisit[0]
		toVisit = toVisit[1:]

		time.Sleep(1 * time.Second)

		normalizedURL := utils.NormalizeURL(url)

		visitLock.Lock()
		if visited[normalizedURL] {
			visitLock.Unlock()
			continue
		}
		visited[normalizedURL] = true
		visitLock.Unlock()

		fmt.Println("Fetching:", url)

		links, err := ScrapeAndExtractLinks(url)
		if err != nil {
			log.Printf("Error scraping %s: %v\n", url, err)
			continue
		}

		//fmt.Printf("Found links on %s: %v\n", url, links)

		for _, link := range links {
			normalizedLink := utils.NormalizeURL(link)
			visitLock.Lock()
			if !visited[normalizedLink] {
				toVisit = append(toVisit, link)
				//fmt.Printf("Added to visit list: %s\n", normalizedLink)
			}
			visitLock.Unlock()
		}
	}

	wg.Wait()
}

func ScrapeAndExtractLinks(pageURL string) ([]string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

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

	if err := utils.SavePageToFile(pageData); err != nil {
		return nil, err
	}

	var links []string
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a[href]')).map(a => a.href)`, &links),
	)
	if err != nil {
		return nil, fmt.Errorf("error extracting links from %s: %v", pageURL, err)
	}

	//fmt.Printf("Links found on %s: %v\n", pageURL, links)

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

	//fmt.Printf("Internal links on %s: %v\n", pageURL, internalLinks)

	return internalLinks, nil
}
