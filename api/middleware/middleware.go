package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/utils"
	"github.com/golang-jwt/jwt/v5"
)

// VerifyJWTMiddleware middleware handles JWT Authorization Bearer verification.
func VerifyJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			errorutils.WriteJSONError(w, "Authorization header missing or invalid access token", http.StatusUnauthorized)
			return
		}
		// Expecting header format: "Bearer <token>".
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			errorutils.WriteJSONError(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		err := utils.VerifyToken(token)

		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				log.Printf("%s", err)
				errorutils.WriteJSONError(w, "token expired", http.StatusForbidden)
				return
			case errors.Is(err, jwt.ErrTokenMalformed):
				log.Printf("%s", err)
				errorutils.WriteJSONError(w, "token malformed", http.StatusForbidden)
				return
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				log.Printf("%s", err)
				errorutils.WriteJSONError(w, "token signature invalid", http.StatusForbidden)
				return
			default:
				log.Printf("%s", err)
				errorutils.WriteJSONError(w, "failed to verify token", http.StatusInternalServerError)
				return
			}
		}
		log.Printf("checked jwt token %s", token)
		next.ServeHTTP(w, r)
	})
}

// LogginMiddleware handles logging endpoint access.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// type AppHandler func(http.ResponseWriter, *http.Request) error

// func WithErrorHandler(h AppHandler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if err := h(w, r); err != nil {
// 			log.Printf("handling %q: %v", r.RequestURI, err)
// 			errorutils.WriteJSONError(w, "something went wrong", http.StatusInternalServerError)
// 		}
// 	})
// }
