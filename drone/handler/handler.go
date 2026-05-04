package handler

import (
	"os"

	"github.com/ahr-i/aero-watch/drone/db"
	"github.com/ahr-i/aero-watch/drone/utils/logging"
	"github.com/gorilla/mux"
)

func CreateHandler() *Handler {
	mux := mux.NewRouter()
	store, err := db.NewStore()
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	if err := store.Init(); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	logging.Info("Successfully initialized database.")

	handler := &Handler{
		Handler: mux,
		store:   store,
	}

	mux.HandleFunc("/ping", handler.pingHandler).Methods("GET")
	mux.HandleFunc("/internal/drone/validate", handler.validateDroneHandler).Methods("POST")
	mux.HandleFunc("/internal/drone", handler.registerDroneHandler).Methods("POST")
	mux.HandleFunc("/internal/drone/status", handler.updateDroneStatusHandler).Methods("PUT")

	return handler
}
