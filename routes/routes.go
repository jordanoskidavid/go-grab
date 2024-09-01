package routes

import (
	"WebScraper/handlers"
	"WebScraper/middleware"
	"net/http"
)

func SetupRoutes() {
	http.Handle("/api/crawl", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.StartCrawlHandler))))
	http.Handle("/api/get-data", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.GetScrapedDataHandler))))
	http.Handle("/api/logout", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.LogoutHandler))))

	http.Handle("//api/delete-data", middleware.JWTAuthMiddleware(middleware.RequireRole("admin", http.HandlerFunc(handlers.DeleteScrapedData))))

	http.HandleFunc("/api/register-user", handlers.RegisterHandler)
	http.HandleFunc("/api/login-user", handlers.LoginHandler)
	//http.HandleFunc("/api/crawl", handlers.StartCrawlHandler)
}
