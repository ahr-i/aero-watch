package handler

import (
	"net/http"

	"github.com/ahr-i/auth/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
	store db.Store
}

type loginRequestBody struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type signupRequestBody struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type signupResponseBody struct {
	User string `json:"user"`
	Role string `json:"role"`
}

type updateRoleRequestBody struct {
	User string `json:"user"`
	Role string `json:"role"`
}

type updateRoleResponseBody struct {
	User string `json:"user"`
	Role string `json:"role"`
}

type loginResponseBody struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type verifyResponseBody struct {
	Valid bool   `json:"valid"`
	User  string `json:"user"`
	Role  string `json:"role"`
}

type roleResponseBody struct {
	Role string `json:"role"`
}

type errorResponseBody struct {
	Error string `json:"error"`
}

type authClaims struct {
	User string `json:"user"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
