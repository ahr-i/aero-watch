package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ahr-i/aero-watch/drone-operation/db"
)

func (h *Handler) validateDroneHandler(w http.ResponseWriter, r *http.Request) {
	var body droneRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.Group == "" || body.Code == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "group and code are required"})
		return
	}

	err := h.store.ValidateDroneModel(body.Group, body.Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rend.JSON(w, http.StatusNotFound, errorResponseBody{Error: "drone model not found"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to validate drone model"})
		return
	}

	rend.JSON(w, http.StatusOK, okayResponseBody{Status: "okay"})
}

func (h *Handler) registerDroneHandler(w http.ResponseWriter, r *http.Request) {
	var body droneRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.Group == "" || body.Code == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "group and code are required"})
		return
	}

	err := h.store.RegisterDroneModel(body.Group, body.Code)
	if err != nil {
		if errors.Is(err, db.ErrDroneAlreadyExists) {
			rend.JSON(w, http.StatusConflict, errorResponseBody{Error: "drone model already exists"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to register drone model"})
		return
	}

	rend.JSON(w, http.StatusCreated, okayResponseBody{Status: "okay"})
}

func (h *Handler) updateDroneStatusHandler(w http.ResponseWriter, r *http.Request) {
	var body droneStatusRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.Group == "" || body.Code == "" || body.Status == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "group, code and status are required"})
		return
	}

	err := h.store.UpdateDroneStatus(body.Group, body.Code, body.Status)
	if err != nil {
		if errors.Is(err, db.ErrInvalidStatus) {
			rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "status must be active or inactive"})
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			rend.JSON(w, http.StatusNotFound, errorResponseBody{Error: "drone model not found"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to update drone status"})
		return
	}

	rend.JSON(w, http.StatusOK, droneStatusResponseBody{
		Status:      "okay",
		DroneStatus: body.Status,
	})
}
