package main

import (
	"net/http"
	"time"

	"github.com/ahr-i/aero-watch/gps-tracking/handler"
	"github.com/ahr-i/aero-watch/gps-tracking/setting"
	"github.com/ahr-i/aero-watch/gps-tracking/utils/logging"

	"github.com/urfave/negroni"
)

func initiallization() {
	logging.Init()
	setting.Init()
}

func startServerHTTP() {
	mux := handler.CreateHandler()
	handler := negroni.Classic()
	defer mux.Close()

	//handler.Use(corsController.SetCors("*", "GET, POST, PUT, DELETE, OPTIONS", "*", true))
	handler.UseHandler(mux)

	logging.Info("HTTP server start.")
	server := &http.Server{
		Addr:              ":" + setting.Setting.ServerPort,
		Handler:           handler,
		ReadHeaderTimeout: time.Duration(setting.Setting.ServerReadHeaderTimeoutSec) * time.Second,
		ReadTimeout:       time.Duration(setting.Setting.ServerReadTimeoutSec) * time.Second,
		WriteTimeout:      time.Duration(setting.Setting.ServerWriteTimeoutSec) * time.Second,
		IdleTimeout:       time.Duration(setting.Setting.ServerIdleTimeoutSec) * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		logging.Error(err)
	}
}

func main() {
	initiallization()

	startServerHTTP()
}
