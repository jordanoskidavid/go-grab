package handlers

import (
	"GoGrab/database"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

// LogoutHandler godoc
// @Summary Logout user
// @Description Logs out the user by deleting the JWT token from the database.
// @Tags Auth
// @Security BearerAuth
// @Produce  plain
// @Success 200 {string} string "Logged out successfully"
// @Failure 401 {string} string "Authorization header missing or invalid token"
// @Failure 500 {string} string "Failed to logout"
// @Router /api/logout [post]

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// check if the request method is GET
	if r.Method != http.MethodGet {
		// If the request method is not GET, return a 405 method not allowed error
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// retrieve the Authorization header from the request
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// if the Authorization header is missing, return a 401 Unauthorized error
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	// extract the token part from the "Bearer" scheme
	tokenString := authHeader[len("Bearer "):]

	// parse the JWT token using the provided key (jwtKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// return the signing key used to verify the token
		return jwtKey, nil
	})

	// if there was an error parsing the token or it's not valid, return a 401 Unauthorized error
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// extract claims from the token, specifically ensuring they are in the expected format
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		// if claims are not in the expected format, return a 401 Unauthorized error
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// extract the user ID from the token's "sub" (subject) claim
	userID, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		// if the user ID cannot be converted to an integer, return a 401 Unauthorized error
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	// call the database function to delete the token for the user
	if err := database.DeleteUserToken(userID); err != nil {
		// if there is an error deleting the token, return a 500 Internal Server Error
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// on successful logout, return a 200 OK status with a message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
