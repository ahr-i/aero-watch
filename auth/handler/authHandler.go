package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ahr-i/aero-watch/auth/db"
	"github.com/ahr-i/aero-watch/auth/setting"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) signupHandler(w http.ResponseWriter, r *http.Request) {
	var body signupRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.User == "" || body.Password == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "user and password are required"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to create password hash"})
		return
	}

	err = h.store.CreateUser(body.User, string(passwordHash), setting.Setting.Role.Unverified)
	if err != nil {
		if errors.Is(err, db.ErrUserAlreadyExists) {
			rend.JSON(w, http.StatusConflict, errorResponseBody{Error: "user already exists"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to create user"})
		return
	}

	rend.JSON(w, http.StatusCreated, signupResponseBody{
		User: body.User,
		Role: setting.Setting.Role.Unverified,
	})
}

func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var body loginRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.User == "" || body.Password == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "user and password are required"})
		return
	}

	userInfo, err := h.store.FindUserAuthInfo(body.User)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rend.JSON(w, http.StatusUnauthorized, errorResponseBody{Error: "invalid user or password"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to find user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.PasswordHash), []byte(body.Password)); err != nil {
		rend.JSON(w, http.StatusUnauthorized, errorResponseBody{Error: "invalid user or password"})
		return
	}

	token, expiresIn, err := createAccessToken(body.User, userInfo.Role)
	if err != nil {
		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to create token"})
		return
	}

	rend.JSON(w, http.StatusOK, loginResponseBody{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	})
}

func (h *Handler) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.ListUsers()
	if err != nil {
		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to list users"})
		return
	}

	responseUsers := make([]userResponseBody, 0, len(users))
	for _, user := range users {
		responseUsers = append(responseUsers, userResponseBody{
			User:      user.User,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		})
	}

	rend.JSON(w, http.StatusOK, listUsersResponseBody{Users: responseUsers})
}

func (h *Handler) verifyHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := parseAccessToken(r)
	if err != nil {
		rend.JSON(w, http.StatusUnauthorized, errorResponseBody{Error: "invalid token"})
		return
	}

	rend.JSON(w, http.StatusOK, verifyResponseBody{
		Valid: true,
		User:  claims.User,
		Role:  claims.Role,
	})
}

func (h *Handler) roleHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := parseAccessToken(r)
	if err != nil {
		rend.JSON(w, http.StatusUnauthorized, errorResponseBody{Error: "invalid token"})
		return
	}

	rend.JSON(w, http.StatusOK, roleResponseBody{Role: claims.Role})
}

func (h *Handler) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	var body updateRoleRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "invalid request body"})
		return
	}

	if body.User == "" || body.Role == "" {
		rend.JSON(w, http.StatusBadRequest, errorResponseBody{Error: "user and role are required"})
		return
	}

	err := h.store.UpdateUserRole(body.User, body.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rend.JSON(w, http.StatusNotFound, errorResponseBody{Error: "user not found"})
			return
		}

		rend.JSON(w, http.StatusInternalServerError, errorResponseBody{Error: "failed to update role"})
		return
	}

	rend.JSON(w, http.StatusOK, updateRoleResponseBody{
		User: body.User,
		Role: body.Role,
	})
}

func createAccessToken(user string, role string) (string, int64, error) {
	expireDuration := time.Duration(setting.Setting.JWT.AccessTokenExpireMin) * time.Minute
	if expireDuration <= 0 {
		expireDuration = time.Hour
	}

	now := time.Now()
	expiresAt := now.Add(expireDuration)
	claims := authClaims{
		User: user,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(setting.Setting.JWT.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(expireDuration.Seconds()), nil
}

func parseAccessToken(r *http.Request) (*authClaims, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return nil, errors.New("missing authorization header")
	}

	tokenString, found := strings.CutPrefix(authorization, "Bearer ")
	if !found || tokenString == "" {
		return nil, errors.New("invalid authorization header")
	}

	token, err := jwt.ParseWithClaims(tokenString, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(setting.Setting.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*authClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
