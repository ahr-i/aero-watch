package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/ahr-i/aero-watch/gps-tracking/setting"
	"github.com/ahr-i/aero-watch/gps-tracking/utils/logging"
	"github.com/gorilla/mux"
)

func (h *Handler) updateDroneLocationHandler(w http.ResponseWriter, r *http.Request) {
	var body gpsRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	body.Group = strings.TrimSpace(body.Group)
	body.Code = strings.TrimSpace(body.Code)

	if errMsg := validateGPSRequest(body); errMsg != "" {
		rend.JSON(w, http.StatusBadRequest, errorResponse{Error: errMsg})
		return
	}

	if !validateDrone(body.Group, body.Code) {
		rend.JSON(w, http.StatusForbidden, errorResponse{Error: "drone validation failed"})
		return
	}

	h.gpsStore.Upsert(body.Group, body.Code, body.Latitude, body.Longitude)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getDroneLocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group := strings.TrimSpace(vars["group"])
	code := strings.TrimSpace(vars["code"])

	if group == "" || code == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponse{Error: "group and code are required"})
		return
	}

	position, ok := h.gpsStore.Get(group, code)
	if !ok {
		rend.JSON(w, http.StatusNotFound, errorResponse{Error: "gps information not found"})
		return
	}

	rend.JSON(w, http.StatusOK, position)
}

func (h *Handler) listDroneLocationHandler(w http.ResponseWriter, r *http.Request) {
	rend.JSON(w, http.StatusOK, gpsListResponse{Drones: h.gpsStore.List()})
}

func validateGPSRequest(body gpsRequestBody) string {
	if body.Group == "" || body.Code == "" {
		return "group and code are required"
	}

	if body.Latitude < -90 || body.Latitude > 90 {
		return "latitude must be between -90 and 90"
	}

	if body.Longitude < -180 || body.Longitude > 180 {
		return "longitude must be between -180 and 180"
	}

	return ""
}

func validateDrone(group string, code string) bool {
	if !setting.Setting.DroneValidateEnabled {
		return true
	}

	if setting.Setting.DroneOperationService == "" || setting.Setting.DroneValidatePath == "" {
		return false
	}

	validateURL, err := url.Parse(strings.TrimRight(setting.Setting.DroneOperationService, "/"))
	if err != nil {
		logging.Error(fmt.Sprintf("drone validation url parse failed: %v", err))
		return false
	}

	validateURL.Path = path.Join(validateURL.Path, setting.Setting.DroneValidatePath)

	body := map[string]string{
		"group": group,
		"code":  code,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		logging.Error(fmt.Sprintf("drone validation marshal failed: %v", err))
		return false
	}

	client := http.Client{
		Timeout: time.Duration(setting.Setting.DroneValidateTimeoutSec) * time.Second,
	}
	resp, err := client.Post(validateURL.String(), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logging.Error(fmt.Sprintf("drone validation failed: %v", err))
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
