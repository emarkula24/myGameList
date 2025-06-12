package handler

import (
	"io"
	"log"
	"net/http"
	"os"

	"example.com/mygamelist/utils"
)

func (h *Handler) Search(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query().Get("query")
	query = utils.ParseSearchQuery(query)
	resp, err := h.Env.API.SearchGames(query)
	if err != nil {
		http.Error(w, "Failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SearchGame(w http.ResponseWriter, req *http.Request) {
	apiKey := os.Getenv("GIANT_BOMB_API_KEY")
	guid := req.URL.Query().Get("guid")
	response, err := http.Get("https://www.giantbomb.com/api/game/" + guid + "/?api_key=" + apiKey + "&format=json")
	if err != nil {
		log.Printf("failed to fetch game data %s:", err)
		http.Error(w, "error searching game with guid", http.StatusBadRequest)
	}
	defer response.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, response.Body); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
