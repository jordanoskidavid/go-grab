package handlers

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func GetScrapedDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	folderPath := "scraping_folder"
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		fileInZip, err := zipWriter.Create(path[len(folderPath)+1:])
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy the file's content to the zip archive
		_, err = io.Copy(fileInZip, file)
		return err
	})

	if err != nil {
		http.Error(w, "Failed to zip folder: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="scraped_data.zip"`)
	w.WriteHeader(http.StatusOK)
}
