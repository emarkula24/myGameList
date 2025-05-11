package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/service"
	"example.com/mygamelist/utils"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(env *utils.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}

		var regReq RegisterRequest
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&regReq); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := service.RegisterUser(env.DB, regReq.Username, regReq.Email, regReq.Password)
		if err != nil {
			switch {
			case errors.Is(err, errorutils.ErrUserExists):
				errorutils.WriteJSONError(w, "User already exists", http.StatusBadRequest)
			default:
				log.Printf("Failed to register user: %s", err)
				http.Error(w, "Error adding user: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		type RegisterResponse struct {
			UserID int64 `json:"user_id"`
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RegisterResponse{UserID: userId})
	}

}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(env *utils.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var loginReq LoginRequest
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&loginReq); err != nil {
			errorutils.WriteJSONError(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		jwtToken, err := service.LoginUser(env.DB, loginReq.Username, loginReq.Password)
		if err != nil {
			switch {
			case errors.Is(err, errorutils.ErrPasswordMatch):
				log.Printf("Failed to login user: %s", err)
				errorutils.WriteJSONError(w, "incorrect password", http.StatusUnauthorized)
			default:
				log.Printf("Failed to login user: %s", err)
				errorutils.WriteJSONError(w, "authentication failed", http.StatusUnauthorized)
			}
			return
		}

		type LoginResponse struct {
			AccessToken string `json:"accessToken"`
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(LoginResponse{AccessToken: jwtToken})
	}
}
