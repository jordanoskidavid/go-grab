package main

import (
	"WebScraper/functions"
	"WebScraper/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {

	functions.CheckDatabaseConnection()
	routes.SetupRoutes()

	fmt.Println("Server started at port:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
