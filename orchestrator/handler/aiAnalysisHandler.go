package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) registerAIAnalysisRoutes(mux *mux.Router) {
	mux.HandleFunc("/user/ai-analysis/ai/question", h.aiAnalysisUserProxyHandler).Methods(http.MethodPost)
	mux.HandleFunc("/user/ai-analysis/ai/question-with-image", h.aiAnalysisUserProxyHandler).Methods(http.MethodPost)
}

func (h *Handler) aiAnalysisUserProxyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.authorizeActiveOrAdmin(w, r) {
		return
	}

	stripProxyPrefix(r, "/user/ai-analysis")
	h.aiAnalysisController.Proxy(w, r)
}
