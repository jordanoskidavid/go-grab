package functions

import (
	"WebScraper/models"
	"encoding/json"
	"fmt"
	"os"
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
