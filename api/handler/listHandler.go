package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/service"
	"github.com/patrickmn/go-cache"
)

type ListHandler struct {
	ListService *service.ListService
	Cache       *cache.Cache
}

func NewListHandler(ls *service.ListService) *ListHandler {
	c := cache.New(10*time.Minute, 15*time.Minute)
	return &ListHandler{ListService: ls, Cache: c}
}

type ListRequest struct {
	GameId   int    `json:"game_id"`
	Status   string `json:"status"`
	UserName string `json:"username"`
}

func (h *ListHandler) InsertToList(w http.ResponseWriter, r *http.Request) {

	var listReq ListRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&listReq); err != nil {
		log.Printf("failed to insert game: %s", err)
		errorutils.WriteJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err := h.ListService.PostGame(listReq.GameId, listReq.UserName, listReq.Status)
	if err != nil {
		log.Printf("failed to insert game: %s", err)
		errorutils.WriteJSONError(w, "failed to add game to list", http.StatusInternalServerError)
		return
	}
	h.Cache.Delete(listReq.UserName)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func (h *ListHandler) UpdateList(w http.ResponseWriter, r *http.Request) {
	var updateReq ListRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&updateReq); err != nil {
		log.Printf("failed to update game: %s", err)
		errorutils.WriteJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err := h.ListService.PutGame(updateReq.GameId, updateReq.UserName, updateReq.Status)
	if err != nil {
		log.Printf("failed to update game: %s", err)
		errorutils.WriteJSONError(w, "failed to update list", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (h *ListHandler) GetList(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")
	// if cachedResp, found := h.Cache.Get(username); found {
	// 	log.Print("used the cache")
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	if _, err := w.Write(cachedResp.([]byte)); err != nil {
	// 		log.Printf("failed to write cached response: %s", err)
	// 	}
	// 	return
	// }

	// Extract 'page' and 'limit' query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 2 // Default to 10 items per page
	}

	response, gameListDb, err := h.ListService.GetGameList(username, page, limit)
	if len(gameListDb) == 0 {
		log.Printf("gamelist is empty: %s", err)
		errorutils.WriteJSONError(w, "gamelist is empty", http.StatusBadRequest)
	}
	if err != nil {
		log.Printf("failed to get gamelist: %s", err)
		errorutils.WriteJSONError(w, "failed to get list", http.StatusInternalServerError)
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

	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&gameJSON)
	if err != nil {
		log.Printf("failed to decode json body: %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	switch gameJSON.StatusCode {
	case 1:
		h.Cache.Set(username, bodyBytes, cache.DefaultExpiration)
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
		log.Printf("Gamebomb API status != 200: %d", gameJSON.StatusCode)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
	}

}
