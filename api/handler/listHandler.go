package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/service"
)

type ListHandler struct {
	ListService *service.ListService
}

func NewListHandler(ls *service.ListService) *ListHandler {
	return &ListHandler{ListService: ls}
}

func (h *ListHandler) AddToList(w http.ResponseWriter, r *http.Request) {

	type ListRequest struct {
		GameId string `json:"game_id"`
		Status string `json:"status"`
		UserId string `json:"user_id"`
	}
	var listReq ListRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&listReq); err != nil {
		log.Printf("Failed to insert game: %s", err)
		errorutils.WriteJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err := h.ListService.AddGame(listReq.GameId, listReq.UserId, listReq.Status)
	if err != nil {
		log.Printf("Failed to insert game: %s", err)
		errorutils.WriteJSONError(w, "Failed to add game to list", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
