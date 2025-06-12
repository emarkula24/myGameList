package handler

import (
	"example.com/mygamelist/interfaces"
	"example.com/mygamelist/utils"
)

type Handler struct {
	Env  *utils.Env
	Repo interfaces.Repository
}

func NewHandler(env *utils.Env, repo interfaces.Repository) *Handler {
	return &Handler{Env: env}
}
