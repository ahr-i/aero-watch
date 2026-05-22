package handler

import (
	"net/http"

	"github.com/ahr-i/aero-watch/drone-operation/db"
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

type matchingRequestBody struct {
	DriverID int64  `json:"driverId"`
	Group    string `json:"group"`
	Code     string `json:"code"`
}

type driverInfoRequestBody struct {
	Content string `json:"content"`
}

type updateDriverInfoRequestBody struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

type deleteDriverInfoRequestBody struct {
	ID int64 `json:"id"`
}

type okayResponseBody struct {
	Status string `json:"status"`
}

type droneStatusResponseBody struct {
	Status      string `json:"status"`
	DroneStatus string `json:"droneStatus"`
}

type droneResponseBody struct {
	Group  string `json:"group"`
	Code   string `json:"code"`
	Status string `json:"status"`
}

type listDroneResponseBody struct {
	Drones []droneResponseBody `json:"drones"`
}

type driverInfoResponseBody struct {
	ID      int64               `json:"id"`
	Content string              `json:"content"`
	Drones  []droneResponseBody `json:"drones,omitempty"`
}

type listDriverInfoResponseBody struct {
	Infos []driverInfoResponseBody `json:"infos"`
}

type errorResponseBody struct {
	Error string `json:"error"`
}
