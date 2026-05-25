package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func New(baseURL string) (*Controller, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	controller := &Controller{
		baseURL: parsedURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		proxy: httputil.NewSingleHostReverseProxy(parsedURL),
	}
	controller.proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	return controller, nil
}

func (c *Controller) Proxy(w http.ResponseWriter, r *http.Request) {
	c.proxy.ServeHTTP(w, r)
}

func (c *Controller) VerifyActiveOrAdmin(r *http.Request) (*VerifyResponse, error) {
	verifyResponse, err := c.Verify(r)
	if err != nil {
		return nil, err
	}
	if verifyResponse.Role != "active" && verifyResponse.Role != "admin" {
		return nil, ErrPermissionDenied
	}

	return verifyResponse, nil
}

func (c *Controller) VerifyAdmin(r *http.Request) (*VerifyResponse, error) {
	verifyResponse, err := c.Verify(r)
	if err != nil {
		return nil, err
	}
	if verifyResponse.Role != "admin" {
		return nil, ErrPermissionDenied
	}

	return verifyResponse, nil
}

func (c *Controller) Verify(r *http.Request) (*VerifyResponse, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return nil, ErrMissingAuthorization
	}

	verifyURL := c.baseURL.ResolveReference(&url.URL{Path: "/verify"})
	request, err := http.NewRequest(http.MethodGet, verifyURL.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", authorization)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, ErrInvalidUser
	}

	var verifyResponse VerifyResponse
	if err := json.NewDecoder(response.Body).Decode(&verifyResponse); err != nil {
		return nil, err
	}
	if !verifyResponse.Valid {
		return nil, ErrInvalidUser
	}

	return &verifyResponse, nil
}
