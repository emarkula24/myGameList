package errorutils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrUserExists = errors.New("user already exists")
)

var (
	ErrPasswordMatch = errors.New("incorrect password")
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func WriteJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(map[string]any{
		"error": message,
		"code":  code,
	})
}
