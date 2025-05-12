package handler

import (
	"io"
	"net/http"
	"os"

	"example.com/mygamelist/utils"
)

func Search(env *utils.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		apiKey := os.Getenv("GIANT_BOMB_API_KEY")
		query := req.URL.Query().Get("query")
		resp, err := http.Get("https://www.giantbomb.com/api/search/?api_key=" + apiKey + "&format=json&query=" + query + "&resources=game")
		if err != nil {
			http.Error(w, "Failed to fetch gamedata", http.StatusInternalServerError)
		}

		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := io.Copy(w, resp.Body); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
}
