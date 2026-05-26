package handler

import (
	"github.com/ahr-i/aero-watch/web-client/setting"
	"github.com/gorilla/mux"
)

func CreateHandler() *Handler {
	mux := mux.NewRouter()
	handler := &Handler{
		Handler: mux,
	}

	mux.HandleFunc("/ping", handler.pingHandler).Methods("GET")
	mux.PathPrefix("/").Handler(spaHandler(setting.Setting.DistPath))

	return handler
}
