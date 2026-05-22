package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

func (h *Handler) createDriverInfoHandler(w http.ResponseWriter, r *http.Request) {
	var body driverInfoRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.Content == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "content is required"})
		return
	}

	info, err := h.store.CreateDriverInfo(body.Content)
	if err != nil {
		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to create driver info"})
		return
	}

	rend.JSON(w, http.StatusCreated, driverInfoResponseBody{
		ID:      info.ID,
		Content: info.Content,
	})
}

func (h *Handler) listDriverInfoHandler(w http.ResponseWriter, r *http.Request) {
	infos, err := h.store.ListDriverInfos()
	if err != nil {
		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to list driver infos"})
		return
	}

	responseInfos := make([]driverInfoResponseBody, 0, len(infos))
	for _, info := range infos {
		responseDrones := make([]droneResponseBody, 0, len(info.Drones))
		for _, drone := range info.Drones {
			responseDrones = append(responseDrones, droneResponseBody{
				Group:  drone.Group,
				Code:   drone.Code,
				Status: drone.Status,
			})
		}

		responseInfos = append(responseInfos, driverInfoResponseBody{
			ID:      info.ID,
			Content: info.Content,
			Drones:  responseDrones,
		})
	}

	rend.JSON(w, http.StatusOK, listDriverInfoResponseBody{Infos: responseInfos})
}

func (h *Handler) updateDriverInfoHandler(w http.ResponseWriter, r *http.Request) {
	var body updateDriverInfoRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.ID <= 0 || body.Content == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "id and content are required"})
		return
	}

	err := h.store.UpdateDriverInfo(body.ID, body.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rend.JSON(w, http.StatusNotFound, errorResponseBody{Error: "driver info not found"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to update driver info"})
		return
	}

	rend.JSON(w, http.StatusOK, driverInfoResponseBody{
		ID:      body.ID,
		Content: body.Content,
	})
}

func (h *Handler) deleteDriverInfoHandler(w http.ResponseWriter, r *http.Request) {
	var body deleteDriverInfoRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.ID <= 0 {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "id is required"})
		return
	}

	err := h.store.DeleteDriverInfo(body.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rend.JSON(w, http.StatusNotFound, errorResponseBody{Error: "driver info not found"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to delete driver info"})
		return
	}

	rend.JSON(w, http.StatusOK, okayResponseBody{Status: "okay"})
}
