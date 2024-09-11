package routes

import (
	_ "WebScraper/docs"
	"WebScraper/handlers"
	"WebScraper/middleware"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() {
	http.Handle("/api/crawl", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.StartCrawlHandler))))
	http.Handle("/api/get-data", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.GetScrapedDataHandler))))
	http.Handle("/api/logout", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.LogoutHandler))))

	http.Handle("/api/delete-data", middleware.JWTAuthMiddleware(middleware.RequireRole("admin", http.HandlerFunc(handlers.DeleteScrapedData))))

	http.HandleFunc("/api/register-user", handlers.RegisterHandler)
	http.HandleFunc("/api/login-user", handlers.LoginHandler)

	http.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}
