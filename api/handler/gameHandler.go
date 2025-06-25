package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/interfaces"
	"example.com/mygamelist/utils"
)

// Defines dependencies for GameHandler struct
type GameHandler struct {
	Gbc interfaces.GiantBombClient
}

// Creates a new instance of GameHandler
func NewGameHandler(gbc interfaces.GiantBombClient) *GameHandler {
	return &GameHandler{Gbc: gbc}
}

// GET /games/?query=string
// Requests GameBomb API for a list of game entries based on a query string and relays the received json to the client.
func (h *GameHandler) Search(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query().Get("query")
	query = utils.ParseSearchQuery(query)
	resp, err := h.Gbc.SearchGames(query)
	if err != nil {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.WriteJSONError(w, "Failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.WriteJSONError(w, "Failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.WriteJSONError(w, "Failed to fetch gamedata", http.StatusInternalServerError)
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
		return
	}

	// status_code received is not the HTTP status code, it is a code that the API sends with the JSON.
	// It is safe to assume HTTP code 404 is not a likely scenario but countermeasures are taken.
	// status_code 1 = success
	// status_code 100 = wrong api key

	switch gameJSON.StatusCode {
	case 1:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(bodyBytes); err != nil {
			log.Printf("failed to write response: %s", err)
			errorutils.WriteJSONError(w, "Failed to fetch gamedata", http.StatusInternalServerError)
			return
		}
	case 100:
		log.Printf("Invalid API key %s", err)
		errorutils.WriteJSONError(w, "Failed to fetch gamedata", http.StatusInternalServerError)
		return
	default:
		log.Printf("Gamebomb API returned an unexpected code: %d", gameJSON.StatusCode)
		errorutils.WriteJSONError(w, "Failed to fetch gamedata", http.StatusInternalServerError)
	}

}

// Requests GameBomb API for the information of a game-entry based on GUID, and relays it to the client.
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
