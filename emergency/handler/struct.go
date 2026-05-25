package handler

import (
	"net/http"

	"github.com/ahr-i/aero-watch/emergency/db"
	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
	db db.Controller
}

type requestBody struct {
	Prompt string `json:"prompt"`
	User   string `json:"user"`
}

type userRequestBody struct {
	User string `json:"user"`
}

type dataFileRequestBody struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type nearestEmergencyRequestBody struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}
