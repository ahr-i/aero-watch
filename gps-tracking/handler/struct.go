package handler

import (
	"net/http"

	"github.com/ahr-i/aero-watch/gps-tracking/gps"
	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
	gpsStore *gps.Store
}

type requestBody struct {
	Prompt string `json:"prompt"`
	User   string `json:"user"`
}

type userRequestBody struct {
	User string `json:"user"`
}

type gpsRequestBody struct {
	Group     string  `json:"group"`
	Code      string  `json:"code"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type gpsListResponse struct {
	Drones []gps.Position `json:"drones"`
}

type errorResponse struct {
	Error string `json:"error"`
}
