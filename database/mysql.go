package database

import (
	"GoGrab/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	log.Println("Connecting to MySQL database...")

	dsn := "root:root@tcp(172.28.0.2:3306)/GrabDB?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database!")
	DB = db
	return db, nil
}

func RegisterUser(username, hashedPassword string) error {
	query := "INSERT INTO Users (username, password) VALUES (?, ?)"
	result := DB.Exec(query, username, hashedPassword)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UsernameExists(username string) (bool, error) {
	var count int64
	result := DB.Table("Users").Where("username = ?", username).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM Users WHERE username = ?"
	result := DB.Raw(query, username).Scan(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func SaveUserToken(userID int, token string, expiration time.Time) error {
	query := "UPDATE Users SET Token = ?, token_expires_at = ? WHERE id = ?"
	result := DB.Exec(query, token, expiration, userID)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CheckUserToken(userID int, tokenString string) (bool, error) {
	var count int

	// Use DB.Raw to run the raw SQL query
	query := "SELECT COUNT(*) FROM Users WHERE ID = ? AND Token = ?"
	result := DB.Raw(query, userID, tokenString).Scan(&count)

	// Check if there's an error in the query execution
	if result.Error != nil {
		return false, result.Error
	}

	// Return true if the count is greater than 0 (i.e., token exists for the user)
	return count > 0, nil
}

func DeleteUserToken(userID int) error {

	query := "UPDATE Users SET Token = NULL WHERE ID = ?"
	result := DB.Exec(query, userID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/*


func UpdateUserPassword(username, newPassword string) error {
    query := "UPDATE users SET password = ? WHERE username = ?"
    result := database.DB.Exec(query, newPassword, username)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func DeleteUserByUsername(username string) error {
    query := "DELETE FROM users WHERE username = ?"
    result := database.DB.Exec(query, username)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

*/
