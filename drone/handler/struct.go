package handler

import (
	"net/http"

	"github.com/ahr-i/aero-watch/drone/db"
	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
	store db.Store
}

type droneRequestBody struct {
	Group string `json:"group"`
	Code  string `json:"code"`
}

type droneStatusRequestBody struct {
	Group  string `json:"group"`
	Code   string `json:"code"`
	Status string `json:"status"`
}

type okayResponseBody struct {
	Status string `json:"status"`
}

type droneStatusResponseBody struct {
	Status      string `json:"status"`
	DroneStatus string `json:"droneStatus"`
}

type errorResponseBody struct {
	Error string `json:"error"`
}
