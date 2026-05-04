package main

import (
	"net/http"

	"github.com/ahr-i/aero-watch/drone/handler"
	"github.com/ahr-i/aero-watch/drone/setting"
	"github.com/ahr-i/aero-watch/drone/utils/corsController"
	"github.com/ahr-i/aero-watch/drone/utils/logging"

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

	handler.Use(corsController.SetCors("*", "GET, POST, PUT, DELETE, OPTIONS", "*", true))
	handler.UseHandler(mux)

	logging.Info("HTTP server start.")
	http.ListenAndServe(":"+setting.Setting.ServerPort, handler)
}

func main() {
	initiallization()

	startServerHTTP()
}
