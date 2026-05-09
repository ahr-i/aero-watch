package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/ahr-i/aero-watch/emergency/datafile"
	"github.com/ahr-i/aero-watch/emergency/setting"
	"github.com/gorilla/mux"
)

func (h *Handler) dataFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := datafile.ListFiles()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	rend.JSON(w, http.StatusOK, files)
}

func (h *Handler) importDataHandler(w http.ResponseWriter, r *http.Request) {
	name, date, err := readNameDate(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	csvData, err := datafile.ReadCSV(name, date)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	tableName := datafile.TableName(name, date)
	err = h.db.ImportCSV(tableName, csvData.Columns, csvData.Rows)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	rend.JSON(w, http.StatusOK, map[string]any{
		"tableName": tableName,
		"rowCount":  len(csvData.Rows),
	})
}

func (h *Handler) dropTableHandler(w http.ResponseWriter, r *http.Request) {
	name, date, err := readNameDate(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	tableName := datafile.TableName(name, date)
	err = h.db.DropTable(tableName)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	rend.JSON(w, http.StatusOK, map[string]string{
		"tableName": tableName,
	})
}

func (h *Handler) dbTablesHandler(w http.ResponseWriter, r *http.Request) {
	tables, err := h.db.ListTables()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	rend.JSON(w, http.StatusOK, tables)
}

func (h *Handler) getTableHandler(w http.ResponseWriter, r *http.Request) {
	name, date, err := readNameDate(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	tableData, err := h.db.GetTable(datafile.TableName(name, date))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	rend.JSON(w, http.StatusOK, tableData)
}

func (h *Handler) nearestEmergencyHandler(w http.ResponseWriter, r *http.Request) {
	latitude, longitude, err := readLatitudeLongitude(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	nearest, err := h.db.FindNearest(latitude, longitude, setting.Setting.NearestEmergencyLimit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	rend.JSON(w, http.StatusOK, nearest)
}

func readNameDate(r *http.Request) (string, string, error) {
	var body dataFileRequestBody
	if r.Body != nil {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&body); err != nil && err != io.EOF {
			return "", "", err
		}
	}

	name := firstNotEmpty(body.Name, r.URL.Query().Get("name"))
	date := firstNotEmpty(body.Date, r.URL.Query().Get("date"))
	name = firstNotEmpty(name, mux.Vars(r)["name"])
	date = firstNotEmpty(date, mux.Vars(r)["date"])

	if err := datafile.ValidateNameDate(name, date); err != nil {
		return "", "", err
	}

	return name, date, nil
}

func firstNotEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func readLatitudeLongitude(r *http.Request) (float64, float64, error) {
	var body nearestEmergencyRequestBody
	if r.Body != nil {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&body); err != nil && err != io.EOF {
			return 0, 0, err
		}
	}

	var latitude *float64
	var longitude *float64

	if body.Latitude != nil {
		latitude = body.Latitude
	}
	if body.Longitude != nil {
		longitude = body.Longitude
	}

	if value := r.URL.Query().Get("latitude"); value != "" {
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, 0, err
		}
		latitude = &parsed
	}

	if value := r.URL.Query().Get("longitude"); value != "" {
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, 0, err
		}
		longitude = &parsed
	}

	if latitude == nil || longitude == nil {
		return 0, 0, errors.New("latitude and longitude are required")
	}
	if *latitude < -90 || *latitude > 90 {
		return 0, 0, errors.New("latitude must be between -90 and 90")
	}
	if *longitude < -180 || *longitude > 180 {
		return 0, 0, errors.New("longitude must be between -180 and 180")
	}

	return *latitude, *longitude, nil
}

func writeError(w http.ResponseWriter, status int, err error) {
	rend.JSON(w, status, map[string]string{
		"error": err.Error(),
	})
}
