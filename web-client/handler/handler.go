package handler

import (
	"github.com/gorilla/mux"
)

func CreateHandler() *Handler {
	mux := mux.NewRouter()
	handler := &Handler{
		Handler: mux,
	}

	mux.HandleFunc("/ping", handler.pingHandler).Methods("GET")
	mux.PathPrefix("/").Handler(spaHandler("dist"))

	return handler
}
