package handler

import (
	aiAnalysisController "github.com/ahr-i/aero-watch/ochestrator/serviceController/aiAnalysis"
	authController "github.com/ahr-i/aero-watch/ochestrator/serviceController/auth"
	droneOperationController "github.com/ahr-i/aero-watch/ochestrator/serviceController/droneOperation"
	emergencyController "github.com/ahr-i/aero-watch/ochestrator/serviceController/emergency"
	gpsTrackingController "github.com/ahr-i/aero-watch/ochestrator/serviceController/gpsTracking"
	streamingController "github.com/ahr-i/aero-watch/ochestrator/serviceController/streaming"
	webClientController "github.com/ahr-i/aero-watch/ochestrator/serviceController/webClient"
	"github.com/ahr-i/aero-watch/ochestrator/setting"
	"github.com/ahr-i/aero-watch/ochestrator/utils/logging"
	"github.com/gorilla/mux"
)

func CreateHandler() *Handler {
	mux := mux.NewRouter()

	aiAnalysisService, err := aiAnalysisController.New(setting.Setting.Services.AIAnalysis.BaseURL)
	if err != nil {
		logging.Error(err)
		panic(err)
	}
	authService, err := authController.New(setting.Setting.Services.Auth.BaseURL)
	if err != nil {
		logging.Error(err)
		panic(err)
	}
	droneOperationService, err := droneOperationController.New(setting.Setting.Services.DroneOperation.BaseURL)
	if err != nil {
		logging.Error(err)
		panic(err)
	}
	emergencyService, err := emergencyController.New(setting.Setting.Services.Emergency.BaseURL)
	if err != nil {
		logging.Error(err)
		panic(err)
	}
	gpsTrackingService, err := gpsTrackingController.New(setting.Setting.Services.GPSTracking.BaseURL)
	if err != nil {
		logging.Error(err)
		panic(err)
	}
	streamingService, err := streamingController.New(setting.Setting.Services.Streaming.BaseURL)
	if err != nil {
		logging.Error(err)
		panic(err)
	}
	webClientService, err := webClientController.New(setting.Setting.Services.WebClient.BaseURL)
	if err != nil {
		logging.Error(err)
		panic(err)
	}

	handler := &Handler{
		Handler:                  mux,
		aiAnalysisController:     aiAnalysisService,
		authController:           authService,
		droneOperationController: droneOperationService,
		emergencyController:      emergencyService,
		gpsTrackingController:    gpsTrackingService,
		streamingController:      streamingService,
		webClientController:      webClientService,
	}

	mux.HandleFunc("/ping", handler.pingHandler).Methods("GET")
	handler.registerAIAnalysisRoutes(mux)
	handler.registerAuthRoutes(mux)
	handler.registerDroneOperationRoutes(mux)
	handler.registerEmergencyRoutes(mux)
	handler.registerGPSTrackingRoutes(mux)
	handler.registerStreamingRoutes(mux)
	handler.registerWebClientRoutes(mux)

	return handler
}
