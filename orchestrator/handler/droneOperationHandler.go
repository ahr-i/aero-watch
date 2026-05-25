package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) registerDroneOperationRoutes(mux *mux.Router) {
	mux.HandleFunc("/admin/drone-operation/internal/drone", h.droneOperationAdminProxyHandler).Methods(http.MethodPost)
	mux.HandleFunc("/admin/drone-operation/internal/drone", h.droneOperationAdminProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/admin/drone-operation/internal/drone", h.droneOperationAdminProxyHandler).Methods(http.MethodDelete)
	mux.HandleFunc("/admin/drone-operation/internal/drone/status", h.droneOperationAdminProxyHandler).Methods(http.MethodPut)

	mux.HandleFunc("/admin/drone-operation/internal/driver", h.droneOperationAdminProxyHandler).Methods(http.MethodPost)
	mux.HandleFunc("/admin/drone-operation/internal/driver", h.droneOperationAdminProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/admin/drone-operation/internal/driver", h.droneOperationAdminProxyHandler).Methods(http.MethodPut)
	mux.HandleFunc("/admin/drone-operation/internal/driver", h.droneOperationAdminProxyHandler).Methods(http.MethodDelete)

	mux.HandleFunc("/admin/drone-operation/internal/matching", h.droneOperationAdminProxyHandler).Methods(http.MethodPost)
	mux.HandleFunc("/admin/drone-operation/internal/matching", h.droneOperationAdminProxyHandler).Methods(http.MethodDelete)
	mux.HandleFunc("/user/drone-operation/internal/matching", h.droneOperationUserProxyHandler).Methods(http.MethodGet)
}

func (h *Handler) droneOperationAdminProxyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.authorizeAdmin(w, r) {
		return
	}

	stripProxyPrefix(r, "/admin/drone-operation")
	h.droneOperationController.Proxy(w, r)
}

func (h *Handler) droneOperationUserProxyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.authorizeActiveOrAdmin(w, r) {
		return
	}

	stripProxyPrefix(r, "/user/drone-operation")
	h.droneOperationController.Proxy(w, r)
}
