package middleware

import (
	"WebScraper/models"
	"context"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "user"

func RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromContext(r.Context())
		if err != nil || user.Role != role {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getUserFromContext(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	if !ok {
		return nil, http.ErrNoLocation
	}
	return user, nil
}

func SetUserInContext(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
