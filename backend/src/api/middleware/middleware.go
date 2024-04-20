package middleware

import (
	"encoding/json"
	"net/http"
	"newsletter/src/api/requestheader"

	// "newsletter/src/pkg/auth"
	newslettererror "newsletter/src/pkg/utils/error"
	"newsletter/src/pkg/utils/logger"
	"newsletter/src/pkg/utils/recover"
)

type contextKey string

const (
	JWTClaimsContextKey contextKey = "JWTClaims"
)

type Middleware struct {
	// AuthService auth.UseCase
	Logs logger.Logger
}

func NewMiddleware(
	// authService auth.UseCase,
	logs logger.Logger,
) *Middleware {
	return &Middleware{
		// AuthService: authService,
		Logs: logs,
	}
}

func (m Middleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover.Recover(m.Logs) {
				w.Header().Set(requestheader.ContentType, requestheader.ApplicationJson)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(newslettererror.NewError(newslettererror.InternalServerError, newslettererror.InternalServerErrorMessage))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
