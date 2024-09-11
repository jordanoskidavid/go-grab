package middleware

import (
	"GoGrab/models"
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const userContextKey contextKey = "user"

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// JWTAuthMiddleware checks the JWT token and sets the user information in the context.
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the JWT token
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Convert Subject claim (user ID) from string to int
		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Extract user information from claims and store it in context
		user := &models.User{
			ID:   userID,      // Setting the ID from the JWT claims
			Role: claims.Role, // Setting the Role from the JWT claims
		}

		// Store the user information in the request context
		ctx := SetUserInContext(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromContext(r.Context())
		if err != nil {
			http.Error(w, "Forbidden: User not found", http.StatusForbidden)
			return
		}
		if user.Role != role {
			http.Error(w, "Forbidden: Insufficient role", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getUserFromContext(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}

func SetUserInContext(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
