package handler

import (
	"net/http"

	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
}

type streamRequestBody struct {
	Group string `json:"group"`
	Code  string `json:"code"`
}

type streamCaptureResponseBody struct {
	Group       string `json:"group"`
	Code        string `json:"code"`
	Status      string `json:"status"`
	ImageBase64 string `json:"imageBase64"`
	ImageType   string `json:"imageType"`
}

type streamResponseBody struct {
	Group     string `json:"group"`
	Code      string `json:"code"`
	Status    string `json:"status"`
	RTMPPath  string `json:"rtmpPath,omitempty"`
	HLSURL    string `json:"hlsUrl"`
	StartedAt string `json:"startedAt,omitempty"`
}

type liveStreamResponseBody struct {
	Streams []streamListItem `json:"streams"`
}

type streamListItem struct {
	Group     string `json:"group"`
	Code      string `json:"code"`
	StreamKey string `json:"streamKey"`
	Status    string `json:"status"`
	StartedAt string `json:"startedAt"`
	HLSURL    string `json:"hlsUrl"`
}
