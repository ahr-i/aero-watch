package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func (h *Handler) hlsStreamHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req streamRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	if !isValidStreamPart(req.Group) || !isValidStreamPart(req.Code) {
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	status := streamStatusOffline
	if _, exists := getStream(req.Group, req.Code); exists {
		status = streamStatusLive
	}

	rend.JSON(w, http.StatusOK, streamResponseBody{
		Group:  req.Group,
		Code:   req.Code,
		Status: status,
		HLSURL: hlsURL(req.Group, req.Code),
	})
}

func (h *Handler) liveStreamsHandler(w http.ResponseWriter, r *http.Request) {
	infos := listStreams()
	items := make([]streamListItem, 0, len(infos))

	for _, info := range infos {
		items = append(items, streamListItem{
			Group:     info.Group,
			Code:      info.Code,
			StreamKey: streamKey(info.Group, info.Code),
			Status:    streamStatusLive,
			StartedAt: info.StartedAt.Format(timeFormatRFC3339),
			HLSURL:    hlsURL(info.Group, info.Code),
		})
	}

	rend.JSON(w, http.StatusOK, liveStreamResponseBody{Streams: items})
}

func (h *Handler) streamStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group := vars["group"]
	code := vars["code"]

	if !isValidStreamPart(group) || !isValidStreamPart(code) {
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	status := streamStatusOffline
	startedAt := ""
	if info, exists := getStream(group, code); exists {
		status = streamStatusLive
		startedAt = info.StartedAt.Format(timeFormatRFC3339)
	}

	rend.JSON(w, http.StatusOK, streamResponseBody{
		Group:     group,
		Code:      code,
		Status:    status,
		RTMPPath:  "/live/" + group + "/" + code,
		HLSURL:    hlsURL(group, code),
		StartedAt: startedAt,
	})
}

func (h *Handler) adminStreamsHandler(w http.ResponseWriter, r *http.Request) {
	infos := listStreams()
	items := make([]adminStreamItem, 0, len(infos))

	for _, info := range infos {
		items = append(items, adminStreamItem{
			Group:     info.Group,
			Code:      info.Code,
			StreamKey: streamKey(info.Group, info.Code),
			StartedAt: info.StartedAt.Format(timeFormatRFC3339),
			HLSURL:    hlsURL(info.Group, info.Code),
			AdminPath: "/admin/streams/" + pathEscape(info.Group) + "/" + pathEscape(info.Code),
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	adminStreamsTemplate.Execute(w, adminStreamsPage{Streams: items})
}

func (h *Handler) adminStreamPlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group := vars["group"]
	code := vars["code"]

	if !isValidStreamPart(group) || !isValidStreamPart(code) {
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	adminPlayerTemplate.Execute(w, adminStreamItem{
		Group:     group,
		Code:      code,
		StreamKey: streamKey(group, code),
		HLSURL:    hlsURL(group, code),
	})
}

func hlsURL(group string, code string) string {
	return "/hls/" + url.PathEscape(group) + "/" + url.PathEscape(code) + "/index.m3u8"
}

func isValidStreamPart(value string) bool {
	if value == "" {
		return false
	}

	for _, r := range value {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '_' || r == '-' {
			continue
		}
		return false
	}

	return true
}

func pathEscape(value string) string {
	return url.PathEscape(value)
}

const timeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"

type adminStreamsPage struct {
	Streams []adminStreamItem
}

type adminStreamItem struct {
	Group     string
	Code      string
	StreamKey string
	StartedAt string
	HLSURL    string
	AdminPath string
}

var adminStreamsTemplate = template.Must(template.ParseFiles("./admin/templates/streams.html"))
var adminPlayerTemplate = template.Must(template.ParseFiles("./admin/templates/player.html"))
