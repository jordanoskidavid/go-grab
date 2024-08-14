package utils

import (
	"WebScraper/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func Read_json_urls(filename string) ([]string, error) {

	takeJSON, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening json file: %v", err)
	}
	defer takeJSON.Close()

	var urlData models.URLDatastruct
	decoder := json.NewDecoder(takeJSON)
	if err := decoder.Decode(&urlData); err != nil {
		return nil, fmt.Errorf("error decoding json: %v", err)
	}

	return urlData.URLs, nil
}

func ReadJson(filename string, data interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening JSON file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	return nil
}

func SavePageToFile(page models.PageData) error {
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
