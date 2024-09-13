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
	// check if the method is DELETE
	if r.Method != http.MethodDelete {
		// If the request method is not DELETE, return a 405 method not allowed error
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// where the scraped data is stored
	folderPath := "./scraping_folder"

	// reading all the files from the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		// if an error occurs while reading directory e.g. folder does not exist, throw internal server error code 500
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	// loop through each file found in the folder
	for _, file := range files {

		// constructing the full file path by joining the folder path with the file name.
		filePath := filepath.Join(folderPath, file.Name())

		// attempt to delete the file
		if err := os.Remove(filePath); err != nil {
			// if there is an error throw error code 500, internal server error along with the name of the file that cannot be deleted
			http.Error(w, fmt.Sprintf("Unable to delete file: %s", file.Name()), http.StatusInternalServerError)
			return
		}
	}
	// write that the response type is plain text
	w.Header().Set("Content-Type", "text/plain")

	// sending status 200 O.K that the files are deleted sucessfully
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "All files deleted successfully")
}
