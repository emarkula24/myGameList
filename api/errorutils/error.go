package errorutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrUserExists    = errors.New("user already exists")
	ErrPasswordMatch = errors.New("incorrect password")
	ErrInvalidToken  = errors.New("invalid token")
)

func WriteJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(map[string]any{
		"error": message,
		"code":  code,
	})
	if err != nil {
		fmt.Printf("failed to encode error message %s", err)
	}
}
