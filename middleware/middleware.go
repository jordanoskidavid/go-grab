package middleware

/*
import (
	"WebScraper/models"
	"context"
	"net/http"
)

// Authorization middleware function
func RequireRole(role string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user, err := getUserFromContext(r.Context())
        if err != nil || user.Role != role {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        // Call the next handler if authorized
        http.HandlerFunc(nextHandler).ServeHTTP(w, r)
    }
}

func getUserFromContext(context context.Context) {
	panic("unimplemented")
}
*/
