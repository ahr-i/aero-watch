package aiAnalysis

import "net/http/httputil"

type Controller struct {
	proxy *httputil.ReverseProxy
}
