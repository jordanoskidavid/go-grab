package routes

import (
	"WebScraper/handlers"
	"WebScraper/middleware"
	"net/http"
)

func SetupRoutes() {
	http.Handle("/api/crawl", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.StartCrawlHandler))))
	http.Handle("/api/get-data", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.GetScrapedDataHandler))))

	http.Handle("//api/delete-data", middleware.JWTAuthMiddleware(middleware.RequireRole("admin", http.HandlerFunc(handlers.DeleteScrapedData))))
	http.Handle("/api/generate-api-key", middleware.JWTAuthMiddleware(middleware.RequireRole("admin", http.HandlerFunc(handlers.GenerateAPIKeyHandler))))

	http.HandleFunc("/api/register-user", handlers.RegisterHandler)
	http.HandleFunc("/api/login-user", handlers.LoginHandler)
}
