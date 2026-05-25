package streaming

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func New(baseURL string) (*Controller, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	controller := &Controller{
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
