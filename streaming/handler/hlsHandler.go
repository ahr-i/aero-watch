package handler

import (
	"net/http"
	"path/filepath"

	"github.com/ahr-i/aero-watch/streaming/setting"
)

func (h *Handler) hlsFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	switch filepath.Ext(r.URL.Path) {
	case ".m3u8":
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	case ".ts":
		w.Header().Set("Content-Type", "video/mp2t")
	}

	http.StripPrefix("/hls/", http.FileServer(http.Dir(setting.Setting.HLSRoot))).ServeHTTP(w, r)
}
