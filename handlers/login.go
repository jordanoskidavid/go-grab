package handlers

import (
	"WebScraper/database"
	"WebScraper/models"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// LoginHandler godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token if the login is successful.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param login body models.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string "access_token"
// @Failure 400 {string} string "Invalid request payload or missing fields"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/login-user [post]

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if credentials.Username == "" || credentials.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByUsername(credentials.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   strconv.Itoa(user.ID),
		},
		Role: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Store the JWT token in the database
	if err := database.SaveUserToken(user.ID, tokenString, expirationTime); err != nil {
		log.Printf("Error saving token to database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send the token to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": tokenString})
}
