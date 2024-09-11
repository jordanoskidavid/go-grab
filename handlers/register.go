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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	exists, err := database.UsernameExists(user.Username)
	if err != nil {
		http.Error(w, "Error checking username", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	err = database.RegisterUser(user.Username, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
