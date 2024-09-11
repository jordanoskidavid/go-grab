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
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	tokenString := authHeader[len("Bearer "):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Get the user ID from the token's subject claim
	userID, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	if err := database.DeleteUserToken(userID); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
