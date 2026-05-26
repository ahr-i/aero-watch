package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) registerAuthRoutes(mux *mux.Router) {
	mux.HandleFunc("/all/auth/login", h.authAllProxyHandler).Methods(http.MethodPost)
	mux.HandleFunc("/all/auth/signup", h.authAllProxyHandler).Methods(http.MethodPost)
	mux.HandleFunc("/user/auth/verify", h.authUserProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/user/auth/role", h.authUserProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/admin/auth/role", h.authAdminProxyHandler).Methods(http.MethodPut)
	mux.HandleFunc("/admin/auth/users", h.authAdminProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/admin/auth/users", h.authAdminProxyHandler).Methods(http.MethodDelete)
}

func (h *Handler) authAllProxyHandler(w http.ResponseWriter, r *http.Request) {
	stripProxyPrefix(r, "/all/auth")
	h.authController.Proxy(w, r)
}

func (h *Handler) authUserProxyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.authorizeActiveOrAdmin(w, r) {
		return
	}

	stripProxyPrefix(r, "/user/auth")
	h.authController.Proxy(w, r)
}

func (h *Handler) authAdminProxyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.authorizeAdmin(w, r) {
		return
	}

	stripProxyPrefix(r, "/admin/auth")
	h.authController.Proxy(w, r)
}
