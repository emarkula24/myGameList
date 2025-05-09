package controllers

import (
	"net/http"

	"example.com/mygamelist/models"
	"example.com/mygamelist/utils"
)

func Register(env *utils.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		username := req.URL.Query().Get("username")
		email := req.URL.Query().Get("email")
		password := req.URL.Query().Get("password")

		_, err := models.AddUser(env.DB, username, email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
