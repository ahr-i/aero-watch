package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) registerEmergencyRoutes(mux *mux.Router) {
	mux.HandleFunc("/admin/emergency/csv", h.emergencyAdminProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/admin/emergency/csv/import", h.emergencyAdminProxyHandler).Methods(http.MethodPost)
	mux.HandleFunc("/admin/emergency/csv/table", h.emergencyAdminProxyHandler).Methods(http.MethodDelete)
	mux.HandleFunc("/admin/emergency/csv/tables", h.emergencyAdminProxyHandler).Methods(http.MethodGet)

	mux.HandleFunc("/user/emergency/csv/table/{name}/{date}", h.emergencyUserProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/user/emergency/emergency/nearest", h.emergencyUserProxyHandler).Methods(http.MethodPost)
}

func (h *Handler) emergencyAdminProxyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.authorizeAdmin(w, r) {
		return
	}

	stripProxyPrefix(r, "/admin/emergency")
	h.emergencyController.Proxy(w, r)
}

func (h *Handler) emergencyUserProxyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.authorizeActiveOrAdmin(w, r) {
		return
	}

	stripProxyPrefix(r, "/user/emergency")
	h.emergencyController.Proxy(w, r)
}
