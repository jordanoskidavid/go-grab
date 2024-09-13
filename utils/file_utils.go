package utils

import (
	"GoGrab/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

/*
	 	Read_json_urls reads a JSON file containing URLs and returns a list of URLs.
	 	It expects the file to match the structure of models.URLDatastruct. Returns an error
		if the file can't be opened or decoded.
*/
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

/*
	 	ReadJson reads a generic JSON file into the provided data structure.
		It decodes the content of the JSON file into the data interface.
		Returns an error if the file can't be opened or decoded.
*/
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

/*
	 	SavePageToFile saves the page data (models.PageData) to a JSON file.
		If the file already exists, it reads the existing content, appends the new page,
		and then writes it back to the file. The file is stored in the "scraping_folder" directory.
		Returns an error if file operations fail.
*/
func SavePageToFile(page models.PageData) error {

	folderPath := "./scraping_folder"

	// Ensure the folder for storing pages exists
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}

	// Get base URL to use as the file name
	baseURL, err := getBaseURL(page.URL)
	if err != nil {
		return err
	}
	fileName := sanitizeFileName(baseURL) + ".json"
	filePath := filepath.Join(folderPath, fileName)

	// Open or create the JSON file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read existing pages if any
	var pages []models.PageData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pages); err != nil && err != io.EOF {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	// Add the new page to the list
	pages = append(pages, page)

	// Truncate the file and write the updated pages back
	file.Truncate(0)
	file.Seek(0, 0)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Indent for better readability
	if err := encoder.Encode(pages); err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	// Get the absolute file path and print a success message
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}
	fmt.Printf("Page data saved to: %s\n", absPath)
	return nil
}

/*
		getBaseURL extracts and returns the scheme and host from a URL.
	 	It is used to create a consistent base URL for file naming.
	 	Returns an error if the URL can't be parsed.
*/
func getBaseURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}
	return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host), nil
}

/*
sanitizeFileName creates a safe file name based on the URL's hostname.
If the hostname can't be parsed, it falls back to "invalid_url" or "default".
*/
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
