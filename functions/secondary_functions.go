package functions

import (
	"GoGrab/database"
	"log"
)

func CheckDatabaseConnection() {
	// Attempt to connect to the database using the Connect function from the database package
	db, err := database.Connect()
	if err != nil {
		// If the connection fails, log a fatal error and exit the program
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	// Retrieve the SQL DB instance from the database connection
	sqlDB, err := db.DB()
	if err != nil {
		// If retrieving the SQL DB instance fails, log a fatal error and exit the program
		log.Fatalf("Failed to get SQL DB instance: %v", err)
	}
	// Ping the database to check if it is reachable and responding
	err = sqlDB.Ping()
	if err != nil {
		// If the ping fails, log a fatal error and exit the program
		log.Fatalf("Failed to ping the database: %v", err)
	}
	// If all checks pass, log a success message
	log.Println("Database connection successful!")
}
