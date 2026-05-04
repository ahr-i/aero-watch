package handler

import (
	"os"

	"github.com/ahr-i/auth/db"
	"github.com/ahr-i/auth/utils/logging"
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

	mux.HandleFunc("/signup", handler.signupHandler).Methods("POST")
	mux.HandleFunc("/login", handler.loginHandler).Methods("POST")

	mux.HandleFunc("/verify", handler.verifyHandler).Methods("GET")
	mux.HandleFunc("/role", handler.roleHandler).Methods("GET")
	mux.HandleFunc("/role", handler.updateRoleHandler).Methods("PUT")

	return handler
}
