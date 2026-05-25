package auth

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	ErrMissingAuthorization = errors.New("missing authorization header")
	ErrInvalidUser          = errors.New("invalid user")
	ErrPermissionDenied     = errors.New("permission denied")
)

type Controller struct {
	baseURL *url.URL
	client  *http.Client
	proxy   *httputil.ReverseProxy
}

type VerifyResponse struct {
	Valid bool   `json:"valid"`
	User  string `json:"user"`
	Role  string `json:"role"`
}
