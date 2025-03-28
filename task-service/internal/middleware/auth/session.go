package auth

import (
	"code-tasks/task-service/internal/api/http/types"
	"code-tasks/task-service/internal/usecases"
	"code-tasks/task-service/utils"
	"context"
	"net/http"
	"strings"
)

var UserIDKey utils.ContextKey = "user_id"

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

			s, err := smanager.GetSessionByID(token)
			if err != nil {
				utils.WriteJSON(w, types.ErrorResponse{Error: "Invalid token"}, http.StatusUnauthorized)
				return
			}
			
			// Добавляем в контекст идентификатор пользователя
			ctx := context.WithValue(r.Context(), UserIDKey, int(s.UserID))
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
