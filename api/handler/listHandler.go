package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/service"
	"example.com/mygamelist/utils"
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
	Status   int    `json:"status"`
	UserName string `json:"username"`
	GameName string `json:"gamename"`
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
	err := h.ListService.PostGame(listReq.GameId, listReq.Status, listReq.UserName, listReq.GameName)
	if err != nil {
		log.Printf("failed to insert game: %s", err)
		errorutils.WriteJSONError(w, "failed to add game to list", http.StatusInternalServerError)
		return
	}
	for k := range h.Cache.Items() {
		if strings.HasPrefix(k, listReq.UserName+",") {
			h.Cache.Delete(k)
		}
	}
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
	err := h.ListService.PutGame(updateReq.GameId, updateReq.Status, updateReq.UserName)
	if err != nil {
		log.Printf("failed to update game: %s", err)
		errorutils.WriteJSONError(w, "failed to update list", http.StatusBadRequest)
		return
	}
	for k := range h.Cache.Items() {
		if strings.HasPrefix(k, updateReq.UserName+",") {
			h.Cache.Delete(k)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

// Handler for GET /list.
//
// URL Parameters:
//   - username (string, required): The username whose game list is being requested.
//   - page (int, optional): The page number for pagination. Defaults to 1 if not provided or invalid.
//   - limit (int, optional): The number of items per page. Defaults to 20 if not provided or invalid.
//
// Description:
// Retrieves a paginated list of game IDs and their play statuses from the database for the specified user.
// Then fetches detailed game data from the GiantBomb API using those IDs, combines the play status into said API response
// and returns the response.
//
// Responses:
//   - 200 OK: Successfully fetched and merged data.
//   - 400 Bad Request: Missing username or empty game list.
//   - 500 Internal Server Error: API failure, data unmarshalling errors, or other unexpected issues.
func (h *ListHandler) GetList(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Printf("query parameter for username or gamename is missing")
		errorutils.WriteJSONError(w, "no query parameter for username or gamename", http.StatusBadRequest)
		return
	}
	// Extract 'page' and 'limit' query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 20 // Default to 20 items per page
	}

	cacheKey := fmt.Sprintf("%s,%d,%d", username, page, limit)

	if cachedResp, found := h.Cache.Get(cacheKey); found {
		log.Print("used the cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(cachedResp.([]byte)); err != nil {
			log.Printf("failed to write cached response: %s", err)
		}
		return
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
	combinedApiResponse, err := utils.CombineGameListJSON(gameListDb, bodyBytes)
	if err != nil {
		log.Printf("failed to inject values in apiResp: %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}
	finalResponse, err := json.Marshal(combinedApiResponse)
	if err != nil {
		log.Printf("failed to marshal combined response: %s", err)
		errorutils.WriteJSONError(w, "failed to getch gamedata", http.StatusInternalServerError)
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

	err = json.NewDecoder(bytes.NewReader(finalResponse)).Decode(&gameJSON)
	if err != nil {
		log.Printf("failed to decode json body: %s", err)
		errorutils.WriteJSONError(w, "failed to fetch gamedata", http.StatusInternalServerError)
		return
	}

	switch gameJSON.StatusCode {
	case 1:
		h.Cache.Set(cacheKey, finalResponse, cache.DefaultExpiration)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(finalResponse); err != nil {
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

func (h *ListHandler) GetListItem(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	gameId := r.URL.Query().Get("gameId")
	if username == "" || gameId == "" {
		log.Printf("query parameter for username or gameId is missing")
		errorutils.WriteJSONError(w, "no query parameter for username or gameId", http.StatusBadRequest)
		return
	}
	gameIdInt, err := strconv.Atoi(gameId)
	if err != nil {
		log.Printf("failed to convert gameId to int")
		errorutils.WriteJSONError(w, "failed to fetch game status from gamelist", http.StatusInternalServerError)
		return
	}
	type GameResponse struct {
		GameData any `json:"gamedata"`
	}
	game := h.ListService.GetGameFromList(username, gameIdInt)
	var response GameResponse
	if game == nil {
		response = GameResponse{GameData: map[string]any{}}
	} else {
		response = GameResponse{GameData: game}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode into json")
		errorutils.WriteJSONError(w, "failed to fetch game status from gamelist", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

}
