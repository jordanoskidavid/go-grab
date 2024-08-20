package database

import (
	"WebScraper/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	log.Println("Connecting to MySQL database...")

	dsn := "root:root@tcp(172.27.0.2:3306)/GrabDB?charset=utf8mb4&parseTime=True&loc=Local"

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
	result := DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

/*
func GetUserByUsername(username string) (*models.User, error) {
    var user models.User
    query := "SELECT * FROM users WHERE username = ?"
    result := database.DB.Raw(query, username).Scan(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}




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
