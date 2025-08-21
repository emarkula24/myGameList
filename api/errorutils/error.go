package errorutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrPasswordMatch      = errors.New("incorrect password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrRefreshTokenExists = errors.New("refreshtoken already in db")
)

// Write sends a specified HTTP response with message.
// If string is empty a generic error message is written.
func Write(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if message == "" {
		message = "something went wrong"
	}
	err := json.NewEncoder(w).Encode(map[string]any{
		"error": message,
	})
	if err != nil {
		fmt.Printf("failed to encode error message %s", err)
	}
}
