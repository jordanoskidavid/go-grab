package handlers

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// GetScrapedDataHandler godoc
// @Summary Download scraped data as a ZIP file
// @Description Retrieves all scraped data from the "scraping_folder" and provides it as a downloadable ZIP file.
// @Tags Scraping
// @Produce application/zip
// @Success 200 {file} file "ZIP file containing scraped data"
// @Failure 405 {string} string "Invalid request method"
// @Failure 500 {string} string "Failed to zip folder"
// @Router /api/get-data [get]

func GetScrapedDataHandler(w http.ResponseWriter, r *http.Request) {
	// check if the request method is GET, as only GET is allowed for downloading the ZIP file
	if r.Method != http.MethodGet {
		// if the method is not GET, return a 405 Method Not Allowed error
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// create a new ZIP writer that will write the ZIP file to the response
	zipWriter := zip.NewWriter(w)

	// ensure that the ZIP writer is closed after the function returns
	defer zipWriter.Close()

	// define the path to the folder where the scraped data is stored
	folderPath := "scraping_folder"

	// walk through all the files in the "scraping_folder"
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {

		// if there's an error accessing a file, return the error to terminate the walk
		if err != nil {
			return err
		}
		// skip any non-regular files (e.g., directories, symlinks)
		if !info.Mode().IsRegular() {
			return nil
		}
		// create a file in the ZIP archive with the same relative path as in the folder
		fileInZip, err := zipWriter.Create(path[len(folderPath)+1:])
		if err != nil {
			// return the error if the file can't be added to the ZIP
			return err
		}

		// open the current file in the folder
		file, err := os.Open(path)
		if err != nil {
			// return the error if the file can't be opened
			return err
		}

		// ensure that the file is closed after copying its contents
		defer file.Close()

		// copy the contents of the file into the ZIP archive
		_, err = io.Copy(fileInZip, file)
		// return any errors that occurred during the copy operation
		return err
	})

	// if there was an error during the folder walk or zipping process, return a 500 error
	if err != nil {
		http.Error(w, "Failed to zip folder: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// set the content type to application/zip for the ZIP file download
	w.Header().Set("Content-Type", "application/zip")
	// set the Content-Disposition header to trigger a file download with the given filename
	w.Header().Set("Content-Disposition", `attachment; filename="scraped_data.zip"`)
	/*// write the HTTP 200 OK status code to indicate success
	w.WriteHeader(http.StatusOK) */
}
