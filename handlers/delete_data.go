package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// DeleteScrapedData godoc
// @Summary Deletes all scraped data
// @Description Deletes all files in the scraping folder where the scraped data is stored.
// @Tags Data
// @Accept  json
// @Produce text/plain
// @Success 200 {string} string "All files deleted successfully"
// @Failure 405 {string} string "Invalid request method"
// @Failure 500 {string} string "Unable to read directory or delete file"
// @Router /api/delete-data [delete]

func DeleteScrapedData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	folderPath := "./scraping_folder"

	files, err := os.ReadDir(folderPath)
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())
		if err := os.Remove(filePath); err != nil {
			http.Error(w, fmt.Sprintf("Unable to delete file: %s", file.Name()), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "All files deleted successfully")

}
