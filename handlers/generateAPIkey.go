package handlers

import (
	"WebScraper/functions"
	"encoding/json"
	"net/http"
)

func GenerateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := functions.GenerateAPIKey()
	if err != nil {
		http.Error(w, "Failed to generate API key", http.StatusInternalServerError)
		return
	}
	response := map[string]string{"apiKey": apiKey}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
