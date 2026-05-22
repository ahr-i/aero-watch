package handler

import (
	"errors"
	"net/http"
	"strings"

	authController "github.com/ahr-i/aero-watch/ochestrator/serviceController/auth"
)

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) authorizeActiveOrAdmin(w http.ResponseWriter, r *http.Request) bool {
	_, err := h.authController.VerifyActiveOrAdmin(r)
	if err == nil {
		return true
	}

	h.writeAuthError(w, err)
	return false
}

func (h *Handler) authorizeAdmin(w http.ResponseWriter, r *http.Request) bool {
	_, err := h.authController.VerifyAdmin(r)
	if err == nil {
		return true
	}

	h.writeAuthError(w, err)
	return false
}

func (h *Handler) writeAuthError(w http.ResponseWriter, err error) {
	statusCode := http.StatusUnauthorized
	if errors.Is(err, authController.ErrPermissionDenied) {
		statusCode = http.StatusForbidden
	}
	if !errors.Is(err, authController.ErrMissingAuthorization) &&
		!errors.Is(err, authController.ErrInvalidUser) &&
		!errors.Is(err, authController.ErrPermissionDenied) {
		statusCode = http.StatusBadGateway
	}

	rend.JSON(w, statusCode, errorResponse{Error: err.Error()})
}

func stripProxyPrefix(r *http.Request, prefix string) {
	r.URL.Path = stripPathPrefix(r.URL.Path, prefix)
	if r.URL.RawPath != "" {
		r.URL.RawPath = stripPathPrefix(r.URL.RawPath, prefix)
	}
}

func stripPathPrefix(path string, prefix string) string {
	strippedPath := strings.TrimPrefix(path, prefix)
	if strippedPath == "" {
		return "/"
	}

	return strippedPath
}
