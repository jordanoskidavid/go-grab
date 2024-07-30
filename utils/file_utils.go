package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

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
