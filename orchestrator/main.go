package main

import (
	"net/http"
	"strings"

	"github.com/ahr-i/aero-watch/ochestrator/handler"
	"github.com/ahr-i/aero-watch/ochestrator/setting"
	"github.com/ahr-i/aero-watch/ochestrator/utils/logging"

	"github.com/urfave/negroni"
)

func initiallization() {
	logging.Init()
	setting.Init()
}

func startServerHTTP() {
	mux := handler.CreateHandler()
	handler := negroni.New(negroni.NewRecovery(), requestLogger(), negroni.NewStatic(http.Dir("public")))
	defer mux.Close()

	//handler.Use(corsController.SetCors("*", "GET, POST, PUT, DELETE, OPTIONS", "*", true))
	handler.UseHandler(mux)

	logging.Info("HTTP server start.")
	http.ListenAndServe(":"+setting.Setting.ServerPort, handler)
}

func main() {
	initiallization()

	startServerHTTP()
}

func requestLogger() negroni.Handler {
	logger := negroni.NewLogger()

	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if strings.HasPrefix(r.URL.Path, "/all/streaming/hls/") {
			next(w, r)
			return
		}

		logger.ServeHTTP(w, r, next)
	})
}
