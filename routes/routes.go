package routes

import (
	_ "GoGrab/docs"
	"GoGrab/handlers"
	"GoGrab/middleware"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() {

	//user avaliable routes
	http.Handle("/api/crawl", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.StartCrawlHandler))))
	http.Handle("/api/get-data", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.GetScrapedDataHandler))))
	http.Handle("/api/logout", middleware.JWTAuthMiddleware(middleware.RequireRole("user", http.HandlerFunc(handlers.LogoutHandler))))

	//admin avaliable routes
	http.Handle("/api/delete-data", middleware.JWTAuthMiddleware(middleware.RequireRole("admin", http.HandlerFunc(handlers.DeleteScrapedData))))

	//public avaliable routes
	http.HandleFunc("/api/register-user", handlers.RegisterHandler)
	http.HandleFunc("/api/login-user", handlers.LoginHandler)
	//http.HandleFunc("/api/get-data", handlers.GetScrapedDataHandler) //testing example

	//documentation routes
	http.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}
