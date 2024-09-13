package handlers

import (
	"GoGrab/functions"
	"GoGrab/models"
	"encoding/json"
	"net/http"
)

// StartCrawlHandler godoc
// @Summary Starts a web crawl process
// @Description Initiates a web scraping process by accepting a list of URLs.
// @Tags Crawling
// @Accept json
// @Produce json
// @Param request body models.URLDatastruct true "List of URLs to crawl"
// @Success 200 {string} string "Crawling completed"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 405 {string} string "Invalid request method"
// @Router /api/crawl [post]

func StartCrawlHandler(w http.ResponseWriter, r *http.Request) {
	// writing response that the crawling process has started
	w.Write([]byte("Crawling started"))

	// check if the request method is POST
	if r.Method != http.MethodPost {
		// If the request method is not POST, return a 405 method not allowed error
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// a variable to store the request payload (list of URLs to crawl)
	var requestData models.URLDatastruct

	// decode the incoming JSON request body into the requestData struct
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		// if there is an error in decoding e.g., invalid JSON, return a 400 bad request error
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// close the request body when the function completes to free resources
	defer r.Body.Close()

	// loop through the list of URLs provided in the request
	for _, url := range requestData.URLs {
		// for each URL, initiate the crawling process using the Crawl function from the functions package
		functions.Crawl(url)
	}

	//This here is optionally, if I want to send status code of 200 ----- w.WriteHeader(http.StatusOK)

	// response that the crawling process has completed
	w.Write([]byte("Crawling completed"))
}
