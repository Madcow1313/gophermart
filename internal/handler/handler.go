package handler

import "github.com/Madcow1313/gophermart/internal/dbconnector"

type Handler struct {
	DBConnector *dbconnector.Connector
}

type RegisterDataJSON struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func New() *Handler {
	return &Handler{}
}
