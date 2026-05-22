package handler

import (
	"net/http"

	aiAnalysisController "github.com/ahr-i/aero-watch/ochestrator/serviceController/aiAnalysis"
	authController "github.com/ahr-i/aero-watch/ochestrator/serviceController/auth"
	droneOperationController "github.com/ahr-i/aero-watch/ochestrator/serviceController/droneOperation"
	emergencyController "github.com/ahr-i/aero-watch/ochestrator/serviceController/emergency"
	gpsTrackingController "github.com/ahr-i/aero-watch/ochestrator/serviceController/gpsTracking"
	streamingController "github.com/ahr-i/aero-watch/ochestrator/serviceController/streaming"
	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
	aiAnalysisController     *aiAnalysisController.Controller
	authController           *authController.Controller
	droneOperationController *droneOperationController.Controller
	emergencyController      *emergencyController.Controller
	gpsTrackingController    *gpsTrackingController.Controller
	streamingController      *streamingController.Controller
}

type requestBody struct {
	Prompt string `json:"prompt"`
	User   string `json:"user"`
}

type userRequestBody struct {
	User string `json:"user"`
}
