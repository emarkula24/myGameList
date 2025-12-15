package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/repository"
	"example.com/mygamelist/utils"
	"github.com/patrickmn/go-cache"
)

type GameService interface {
	SearchGames(ctx context.Context, query string) (*http.Response, error)
	SearchGame(ctx context.Context, guid string) (*http.Response, error)
	SearchGameList(ctx context.Context, games []repository.Game, limit int) (*http.Response, error)
}

// Defines a game HTTP handler.
type GameHandler struct {
	Gbc   GameService
	Cache *cache.Cache
}

// NewGameHandler creates a new game HTTP handler.
// The handler uses go-cache to store api responses because the amount of API requests per hour is limited.
func NewGameHandler(gbc GameService) *GameHandler {
	c := cache.New(10*time.Minute, 15*time.Minute)
	return &GameHandler{Gbc: gbc, Cache: c}
}

// Search handles GET /games/search requests.
// Requests GameBomb API for a list of game entries based on a query string and relays the received json to the client.
func (h *GameHandler) Search(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
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
	resp, err := h.Gbc.SearchGames(ctx, query)
	if err != nil {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to fetch gamedata: %s", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}

	err = resp.Body.Close()
	if err != nil {
		errorutils.Write(w, "", http.StatusInternalServerError)
	}

	h.Cache.Set(query, bodyBytes, cache.DefaultExpiration)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bodyBytes); err != nil {
		log.Printf("failed to write response: %s", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}
}

// SearchGame handles GET /games/game requests.
// Requests GameBomb API for the information of a game-entry based on GUID, and relays it to the client.
// Uses go-cache to store response data due to the request amount to gamebomb API being limited
func (h *GameHandler) SearchGame(w http.ResponseWriter, req *http.Request) {
	guid := req.URL.Query().Get("guid")
	ctx := context.Background()
	if cachedResp, found := h.Cache.Get(guid); found {
		log.Print("used the cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(cachedResp.([]byte)); err != nil {
			log.Printf("failed to write cached response: %s", err)
		}
		return
	}

	response, err := h.Gbc.SearchGame(ctx, guid)
	if err != nil {
		log.Printf("failed to fetch game data %s:", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("failed to read response body: %s", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}
	err = response.Body.Close()
	if err != nil {
		log.Printf("failed to close body: %s", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}

	h.Cache.Set(guid, bodyBytes, cache.DefaultExpiration)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bodyBytes); err != nil {
		log.Printf("failed to write response: %s", err)
		errorutils.Write(w, "", http.StatusInternalServerError)
		return
	}
}
