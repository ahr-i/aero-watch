package handler

import (
	"net/http"

	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
}

type requestBody struct {
	Prompt string `json:"prompt"`
}

type imageRequestBody struct {
	Prompt      string `json:"prompt"`
	ImageBase64 string `json:"image_base64"`
	ImageType   string `json:"imageType,omitempty"`
}

type answerResponseBody struct {
	Answer string `json:"answer"`
}

type errorResponseBody struct {
	Error string `json:"error"`
}
