package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

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
