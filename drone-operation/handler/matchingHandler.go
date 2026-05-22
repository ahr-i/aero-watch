package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ahr-i/aero-watch/drone-operation/db"
)

func (h *Handler) createMatchingHandler(w http.ResponseWriter, r *http.Request) {
	var body matchingRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.DriverID <= 0 || body.Group == "" || body.Code == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "driverId, group and code are required"})
		return
	}

	err := h.store.CreateDroneDriverMatch(body.DriverID, body.Group, body.Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rend.JSON(w, http.StatusNotFound, errorResponseBody{Error: "driver or drone not found"})
			return
		}

		if errors.Is(err, db.ErrMatchingAlreadyExists) {
			rend.JSON(w, http.StatusConflict, errorResponseBody{Error: "matching already exists"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to create matching"})
		return
	}

	rend.JSON(w, http.StatusCreated, okayResponseBody{Status: "okay"})
}

func (h *Handler) findMatchingHandler(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	code := r.URL.Query().Get("code")
	if group == "" || code == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "group and code are required"})
		return
	}

	info, err := h.store.FindDriverInfoByDrone(group, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rend.JSON(w, http.StatusNotFound, errorResponseBody{Error: "matching not found"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to find matching"})
		return
	}

	rend.JSON(w, http.StatusOK, driverInfoResponseBody{
		ID:      info.ID,
		Content: info.Content,
	})
}

func (h *Handler) deleteMatchingHandler(w http.ResponseWriter, r *http.Request) {
	var body matchingRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.DriverID <= 0 || body.Group == "" || body.Code == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "driverId, group and code are required"})
		return
	}

	err := h.store.DeleteDroneDriverMatch(body.DriverID, body.Group, body.Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rend.JSON(w, http.StatusNotFound, errorResponseBody{Error: "matching not found"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to delete matching"})
		return
	}

	rend.JSON(w, http.StatusOK, okayResponseBody{Status: "okay"})
}
