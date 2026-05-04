package handler

import (
	"time"

	"github.com/ahr-i/aero-watch/gps-tracking/gps"
	"github.com/ahr-i/aero-watch/gps-tracking/setting"
	"github.com/gorilla/mux"
)

func CreateHandler() *Handler {
	mux := mux.NewRouter()
	handler := &Handler{
		Handler:  mux,
		gpsStore: gps.NewStore(time.Duration(setting.Setting.GPSAliveTimeoutSec)*time.Second, time.Duration(setting.Setting.GPSCleanupIntervalSec)*time.Second),
	}

	mux.HandleFunc("/ping", handler.pingHandler).Methods("GET")
	mux.HandleFunc("/drone/location", handler.updateDroneLocationHandler).Methods("POST")
	mux.HandleFunc("/drone/location", handler.listDroneLocationHandler).Methods("GET")
	mux.HandleFunc("/drone/location/{group}/{code}", handler.getDroneLocationHandler).Methods("GET")

	return handler
}
