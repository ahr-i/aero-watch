package handler

import (
	"github.com/gorilla/mux"
)

func CreateHandler() *Handler {
	mux := mux.NewRouter()
	handler := &Handler{
		Handler: mux,
	}

	mux.HandleFunc("/ping", handler.pingHandler).Methods("GET")
	mux.HandleFunc("/ai/question", handler.questionHandler).Methods("POST")
	mux.HandleFunc("/ai/question-with-image", handler.questionWithImageHandler).Methods("POST")

	return handler
}
