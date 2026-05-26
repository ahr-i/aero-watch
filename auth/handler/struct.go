package handler

import (
	"net/http"
	"time"

	"github.com/ahr-i/aero-watch/auth/db"
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

type deleteUserRequestBody struct {
	User string `json:"user"`
}

type deleteUserResponseBody struct {
	User string `json:"user"`
}

type listUsersResponseBody struct {
	Users []userResponseBody `json:"users"`
}

type userResponseBody struct {
	User      string    `json:"user"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type loginResponseBody struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int64  `json:"expiresIn"`
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
