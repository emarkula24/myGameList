package handler

import (
	"net/http"

	"example.com/mygamelist/service"
	"example.com/mygamelist/utils"
	"github.com/gorilla/mux"
)

func Register(env *utils.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		// username := req.URL.Query().Get("username")
		// email := req.URL.Query().Get("email")
		// password := req.URL.Query().Get("password")

		_, err := service.AddUser(env.DB, vars["username"], vars["email"], vars["password"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
