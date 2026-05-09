package main

import (
	"net/http"
	"os"

	"github.com/ahr-i/aero-watch/emergency/datafile"
	"github.com/ahr-i/aero-watch/emergency/db"
	"github.com/ahr-i/aero-watch/emergency/handler"
	"github.com/ahr-i/aero-watch/emergency/setting"
	"github.com/ahr-i/aero-watch/emergency/utils/corsController"
	"github.com/ahr-i/aero-watch/emergency/utils/logging"

	"github.com/urfave/negroni"
)

func initiallization() {
	logging.Init()
	setting.Init()

	err := datafile.ValidateFiles()
	if err != nil {
		logging.Error(err)

		os.Exit(1)
	}
}

func startServerHTTP() {
	dbController, err := db.NewMySQLController()
	if err != nil {
		logging.Error(err)

		os.Exit(1)
	}

	mux := handler.CreateHandler(dbController)
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
