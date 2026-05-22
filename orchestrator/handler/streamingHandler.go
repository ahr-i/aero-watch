package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) registerStreamingRoutes(mux *mux.Router) {
	mux.HandleFunc("/user/streaming/api/v1/streams/hls", h.streamingUserProxyHandler).Methods(http.MethodPost)
	mux.HandleFunc("/user/streaming/api/v1/streams/live", h.streamingUserProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/user/streaming/api/v1/streams/capture", h.streamingUserProxyHandler).Methods(http.MethodPost)
	mux.HandleFunc("/user/streaming/api/v1/streams/{group}/{code}", h.streamingUserProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/all/streaming/hls/{group}/{code}/index.m3u8", h.streamingAllProxyHandler).Methods(http.MethodGet)
	mux.HandleFunc("/all/streaming/hls/{group}/{code}/{segment}.ts", h.streamingAllProxyHandler).Methods(http.MethodGet)
}

func (h *Handler) streamingUserProxyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.authorizeActiveOrAdmin(w, r) {
		return
	}

	stripProxyPrefix(r, "/user/streaming")
	h.streamingController.Proxy(w, r)
}

func (h *Handler) streamingAllProxyHandler(w http.ResponseWriter, r *http.Request) {
	stripProxyPrefix(r, "/all/streaming")
	h.streamingController.Proxy(w, r)
}
