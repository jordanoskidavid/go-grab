package handlers

import (
	"GoGrab/database"
	"GoGrab/models"
	"GoGrab/utils"
	"encoding/json"
	"net/http"
)

// RegisterHandler godoc
// @Summary Register a new user
// @Description Registers a new user with a username and password. The password is hashed before saving.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param  user  body  models.User  true  "User data"
// @Success 201 {object} map[string]string "User registered successfully"
// @Failure 400 {string} string "Username and password are required or Invalid request payload"
// @Failure 409 {string} string "Username already exists"
// @Failure 500 {string} string "Failed to register user or Error checking username"
// @Router /api/register-user [post]

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	// check if the request method is POST
	if r.Method != http.MethodPost {
		// If the request method is not POST, return a 405 method not allowed error
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// declare a variable to hold the user data
	var user models.User

	// create a JSON decoder to parse the request body into the user struct
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		// if the request body is not valid JSON or fails to decode, return a 400 Bad Request error
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// check if the username already exists in the database
	exists, err := database.UsernameExists(user.Username)
	if err != nil {
		// if there is an error checking the username, return a 500 Internal Server Error
		http.Error(w, "Error checking username", http.StatusInternalServerError)
		return
	}

	// if the username already exists, return a 409 Conflict error
	if exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// ensure both username and password are provided, if not, return a 400 Bad Request error
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// hash the user's password before storing it in the database
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		// if there is an error hashing the password, return a 500 Internal Server Error
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// register the user by saving the username and hashed password in the database
	err = database.RegisterUser(user.Username, hashedPassword)
	if err != nil {
		// if registration fails, return a 500 Internal Server Error
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// On successful registration, return a 201 Created status with a success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
