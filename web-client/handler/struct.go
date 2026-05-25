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
	User   string `json:"user"`
}

type userRequestBody struct {
	User string `json:"user"`
}
