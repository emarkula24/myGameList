package handler

import (
	"encoding/json"
	"net/http"

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

		_, err := service.RegisterUser(env.DB, regReq.Username, regReq.Email, regReq.Password)
		if err != nil {
			http.Error(w, "Error adding user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		jwtToken, err := service.LoginUser(env.DB, loginReq.Username, loginReq.Password)
		if err != nil {
			http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"accessToken": jwtToken})
	}
}
