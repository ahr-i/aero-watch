package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) registerWebClientRoutes(mux *mux.Router) {
	mux.PathPrefix("/").HandlerFunc(h.webClientProxyHandler)
}

func (h *Handler) webClientProxyHandler(w http.ResponseWriter, r *http.Request) {
	h.webClientController.Proxy(w, r)
}
