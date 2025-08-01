package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/interfaces"
	"example.com/mygamelist/utils"
	"github.com/patrickmn/go-cache"
)

// Defines dependencies for GameHandler struct
type GameHandler struct {
	Gbc   interfaces.GiantBombClient
	Cache *cache.Cache
}

// Creates a new instance of GameHandler.
// The handler uses go-cache to store api responses because the amount of API requests per hour is limited.
func NewGameHandler(gbc interfaces.GiantBombClient) *GameHandler {
	c := cache.New(10*time.Minute, 15*time.Minute)
	return &GameHandler{Gbc: gbc, Cache: c}
}

// GET /games/?query=string
// Requests GameBomb API for a list of game entries based on a query string and relays the received json to the client.
func (h *GameHandler) Search(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query().Get("query")
	query = utils.ParseSearchQuery(query)

	if cachedResp, found := h.Cache.Get(query); found {
		log.Print("used the cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(cachedResp.([]byte)); err != nil {
			log.Printf("failed to write cached response: %s", err)
		}
		return
	}
	resp, err := h.Gbc.SearchGames(query)
	if err != nil {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
	}

	type GameJSON struct {
		StatusCode int `json:"status_code"`
	}
	var gameJSON GameJSON

	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&gameJSON)
	if err != nil {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	// status_code received is not the HTTP status code, it is a code that the API sends with the JSON.
	// It is safe to assume HTTP code 404 is not a likely scenario but countermeasures are taken.
	// status_code 1 = success
	// status_code 100 = wrong api key

	switch gameJSON.StatusCode {
	case 1:
		h.Cache.Set(query, bodyBytes, cache.DefaultExpiration)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(bodyBytes); err != nil {
			log.Printf("failed to write response: %s", err)
			errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
			return
		}
	case 100:
		log.Printf("Invalid API key %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	default:
		log.Printf("Gamebomb API returned an unexpected code: %d", gameJSON.StatusCode)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
	}

}

// Requests GameBomb API for the information of a game-entry based on GUID, and relays it to the client.
// Uses go-cache to store response data due to the request amount to gamebomb API being limited
func (h *GameHandler) SearchGame(w http.ResponseWriter, req *http.Request) {
	guid := req.URL.Query().Get("guid")

	if cachedResp, found := h.Cache.Get(guid); found {
		log.Print("used the cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(cachedResp.([]byte)); err != nil {
			log.Printf("failed to write cached response: %s", err)
		}
		return
	}

	response, err := h.Gbc.SearchGame(guid)
	if err != nil {
		log.Printf("failed to fetch game data %s:", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("failed to read response body: %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}
	err = response.Body.Close()
	if err != nil {
		log.Printf("failed to close body: %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}
	type GameJSON struct {
		StatusCode int `json:"status_code"`
	}
	var gameJSON GameJSON

	// status_code received is not the HTTP status code, it is a code that the API sends with the JSON.
	// It is safe to assume HTTP code 404 is not a likely scenario but countermeasures are taken.
	// status_code 1 = success
	// status_code 100 = wrong api key
	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&gameJSON)
	if err != nil {
		log.Printf("failed to decode json body: %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	switch gameJSON.StatusCode {
	case 1:
		h.Cache.Set(guid, bodyBytes, cache.DefaultExpiration)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(bodyBytes); err != nil {
			log.Printf("failed to write response: %s", err)
			errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
			return
		}
	case 100:
		log.Printf("Invalid API key %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	default:
		log.Printf("Gamebomb API returned an unexpected code: %d", gameJSON.StatusCode)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
	}
}
