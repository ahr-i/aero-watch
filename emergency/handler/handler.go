package handler

import (
	"github.com/ahr-i/aero-watch/emergency/db"
	"github.com/gorilla/mux"
)

func CreateHandler(dbController db.Controller) *Handler {
	mux := mux.NewRouter()
	handler := &Handler{
		Handler: mux,
		db:      dbController,
	}

	mux.HandleFunc("/ping", handler.pingHandler).Methods("GET")
	mux.HandleFunc("/csv", handler.dataFilesHandler).Methods("GET")
	mux.HandleFunc("/csv/import", handler.importDataHandler).Methods("POST")
	mux.HandleFunc("/csv/table", handler.dropTableHandler).Methods("DELETE")
	mux.HandleFunc("/csv/tables", handler.dbTablesHandler).Methods("GET")
	mux.HandleFunc("/csv/table/{name}/{date}", handler.getTableHandler).Methods("GET")
	mux.HandleFunc("/emergency/nearest", handler.nearestEmergencyHandler).Methods("POST")

	return handler
}
