package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/utils"
)

func (h *GameHandler) Search(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query().Get("query")
	query = utils.ParseSearchQuery(query)
	resp, err := h.Env.API.SearchGames(query)
	if err != nil {
		http.Error(w, "Failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	type GameJSON struct {
		StatusCode int `json:"status_code"`
	}
	var gameJSON GameJSON

	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&gameJSON)
	if err != nil {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.WriteJSONError(w, "Failed to fetch gamedata", http.StatusInternalServerError)
	}

	if gameJSON.StatusCode != 1 {
		log.Printf("data fetching from gamebomb failed")
		http.Error(w, "Failed to fetch gamedata", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(bodyBytes); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) SearchGame(w http.ResponseWriter, req *http.Request) {
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
