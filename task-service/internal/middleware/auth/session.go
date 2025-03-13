package auth

import (
	"net/http"
	"strings"
	"task-service/internal/api/http/types"
	"task-service/internal/usecases"
	"task-service/utils"
)

var authPrefix = "Bearer "

func SessionMiddleware(smanager usecases.Session) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteJSON(w, types.ErrorResponse{Error: "Missing token"}, http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, authPrefix) {
				utils.WriteJSON(w, types.ErrorResponse{Error: "Invalid token format"}, http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, authPrefix)

			_, err := smanager.GetSessionByID(token)
			if err != nil {
				utils.WriteJSON(w, types.ErrorResponse{Error: "Invalid token"}, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
