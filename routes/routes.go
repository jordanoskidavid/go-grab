package routes

import (
	"WebScraper/handlers"
	"WebScraper/middleware"
	"net/http"
)

func SetupRoutes() {
	//protected routes down here

	http.Handle("/api/crawl", middleware.RequireRole("user", http.HandlerFunc(handlers.StartCrawlHandler)))
	http.Handle("/api/get-data", middleware.RequireRole("user", http.HandlerFunc(handlers.GetScrapedDataHandler)))

	http.Handle("/api/delete-data", middleware.RequireRole("admin", http.HandlerFunc(handlers.DeleteScrapedData)))
	http.Handle("/api/generate-api-key", middleware.RequireRole("admin", http.HandlerFunc(handlers.GenerateAPIKeyHandler)))

	//public routes down here

	http.HandleFunc("/api/register-user", handlers.RegisterHandler)
	http.HandleFunc("/api/login-user", handlers.LoginHandler)
}
