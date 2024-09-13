package handlers

import (
	"GoGrab/database"
	"GoGrab/models"
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

	// check if the request method is POST
	if r.Method != http.MethodPost {
		// If the request method is not POST, return a 405 method not allowed error
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// temp struct for storing login credentials
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Decoding the json body from the request into the credentials struct
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		// if the json is invalid, it returns error code 400 for bad request
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// here checks if the username or password is empty
	if credentials.Username == "" || credentials.Password == "" {
		// if it's empty, returns bad requesst error code 400
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// retrieving the user from the database by username
	user, err := database.GetUserByUsername(credentials.Username)
	if err != nil {
		// if the username is invalid it returns error code 401, unauthorized
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// here compares provided password with the hashed password stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		// and again here, if the password doesn't match, it returns error code 401, unauthorized
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// here, I set the token expiration time to 60 minutes from now
	expirationTime := time.Now().Add(60 * time.Minute)
	// I use JWT claims, i.e creating for getting the metadata for user id as subject and role,  for not goint repeatedly in the database and checking
	claims := &models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   strconv.Itoa(user.ID),
		},
		Role: user.Role,
	}
	// creating the jwt token using HS256 signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signing the token using the secret key, in my case it's "JWT_SECRET_KEY" and because is not used in production
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// if it fails returns code error 500, internal server error
		log.Printf("Error signing token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// save the generated JWT token to the database for the user along with the expiration time
	if err := database.SaveUserToken(user.ID, tokenString, expirationTime); err != nil {
		// if saving the token to the database fails, will return error 500, internal server error
		log.Printf("Error saving token to database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// setting the response header to indicate json content
	w.Header().Set("Content-Type", "application/json")

	// return the signed jwt token as a json response
	json.NewEncoder(w).Encode(map[string]string{"access_token": tokenString})
}
