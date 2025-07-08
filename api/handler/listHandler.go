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

type ListRequest struct {
	GameId int    `json:"game_id"`
	Status string `json:"status"`
	UserId int    `json:"user_id"`
}

func (h *ListHandler) AddToList(w http.ResponseWriter, r *http.Request) {

	var listReq ListRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&listReq); err != nil {
		log.Printf("failed to insert game: %s", err)
		errorutils.WriteJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err := h.ListService.PostGame(listReq.GameId, listReq.UserId, listReq.Status)
	if err != nil {
		log.Printf("failed to insert game: %s", err)
		errorutils.WriteJSONError(w, "failed to add game to list", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

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
	err := h.ListService.PutGame(updateReq.GameId, updateReq.UserId, updateReq.Status)
	if err != nil {
		log.Printf("failed to update game: %s", err)
		errorutils.WriteJSONError(w, "failed to update list", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
