package handler

import (
	"example.com/mygamelist/interfaces"
	"example.com/mygamelist/utils"
)

type Handler struct {
	Env  *utils.Env
	Repo interfaces.Repository
}

type GameHandler struct {
	Env *utils.Env
}

func NewGameHandler(env *utils.Env) *GameHandler {
	return &GameHandler{Env: env}
}
func NewHandler(env *utils.Env, repo interfaces.Repository) *Handler {
	return &Handler{Env: env}
}
