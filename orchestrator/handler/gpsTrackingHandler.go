package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) registerGPSTrackingRoutes(mux *mux.Router) {
	mux.HandleFunc("/all/gps-tracking/drone/location", h.gpsTrackingAllProxyHandler).Methods(http.MethodPost)
	mux.HandleFunc("/all/gps-tracking/drone/location", h.gpsTrackingAllProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/all/gps-tracking/drone/location/{group}/{code}", h.gpsTrackingAllProxyHandler).Methods(http.MethodGet)
}

func (h *Handler) gpsTrackingAllProxyHandler(w http.ResponseWriter, r *http.Request) {
	stripProxyPrefix(r, "/all/gps-tracking")
	h.gpsTrackingController.Proxy(w, r)
}

func (h *Handler) gpsTrackingUserProxyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.authorizeActiveOrAdmin(w, r) {
		return
	}

	stripProxyPrefix(r, "/user/gps-tracking")
	h.gpsTrackingController.Proxy(w, r)
}
