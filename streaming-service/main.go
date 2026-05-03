package main

import (
	"net/http"

	"github.com/ahr-i/aero-watch/handler"
	"github.com/ahr-i/aero-watch/setting"
	"github.com/ahr-i/aero-watch/utils/corsController"
	"github.com/ahr-i/aero-watch/utils/logging"

	"github.com/urfave/negroni"
)

func initiallization() {
	logging.Init()
	setting.Init()
}

func startServerHTTP() {
	mux := handler.CreateHandler()
	handler := negroni.New()
	defer mux.Close()

	handler.Use(negroni.NewRecovery())
	handler.Use(corsController.SetCors("*", "GET, POST, PUT, DELETE, OPTIONS", "*", true))
	handler.UseHandler(mux)

	logging.Info("HTTP server start.")
	http.ListenAndServe(":"+setting.Setting.ServerPort, handler)
}

func main() {
	initiallization()

	go handler.StartRTMPServer()
	startServerHTTP()
}
