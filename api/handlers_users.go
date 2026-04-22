package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabeamv/overloaded-api/internal/auth"
	"github.com/gabeamv/overloaded-api/internal/database"
	"github.com/google/uuid"
)

const (
	SECOND = 1
	MINUTE = SECOND * 60
	HOUR   = MINUTE * 60
	DAY    = HOUR * 24
)

type respUser struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Config) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&req)
	if err != nil {
		err = fmt.Errorf("error decoding for %v: %w", r.Body, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	now := time.Now().UTC()
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		err = fmt.Errorf("error hashing password from request '%v': %w", req, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	hashedPasswordNullStr := sql.NullString{String: hashedPassword, Valid: true}
	user, err := c.DbQueries.CreateUser(context.Background(), database.CreateUserParams{CreatedAt: now, UpdatedAt: now, Email: req.Email,
		HashedPassword: hashedPasswordNullStr, Username: req.Username})
	if err != nil {
		err = fmt.Errorf("error creating user. req=%v: %w", req, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	resp := respUser{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, Username: req.Username}
	ResponseJSON(w, http.StatusCreated, resp)
}

func (c *Config) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Id           uuid.UUID `json:"id"`
		Username     string    `json:"username"`
		Email        string    `json:"email"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&req)
	if err != nil {
		err = fmt.Errorf("error decoding '%v': %w", r.Body, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	user, err := c.DbQueries.GetUserByUsername(context.Background(), req.Username)
	if err != nil {
		err = fmt.Errorf("error getting user using email '%v': %w", req.Username, err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	match, err := auth.CheckPasswordHash(req.Password, user.HashedPassword.String)
	if err != nil {
		err = fmt.Errorf("error checking password hash for '%v': %w", req.Username, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if !match {
		err := fmt.Errorf("error. no password match for email '%v': %w", req.Username, err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	accessExpireInSeconds := HOUR
	accessToken, err := auth.MakeJWT(user.ID, c.Secret, time.Duration(accessExpireInSeconds)*time.Second)
	if err != nil {
		err = fmt.Errorf("error making jwt for user '%v': %w", user, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		err = fmt.Errorf("error generating refresh token for user '%v': %w", user, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	now := time.Now().UTC()
	refreshExpireInSeconds := DAY * 60
	refreshExpire := now.Add(time.Duration(refreshExpireInSeconds) * time.Second)
	refreshToken, err := c.DbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		ExpireAt:  refreshExpire,
	})
	if err != nil {
		err = fmt.Errorf("error adding refresh token to database for user '%v': %w", user.ID, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	resp := response{Id: user.ID, Username: req.Username, Email: user.Email, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
		Token: accessToken, RefreshToken: refreshToken.Token}
	ResponseJSON(w, http.StatusAccepted, resp)
}

func (c *Config) HandlerUpdateUserUsernameById(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Username string `json:"username"`
	}
	userId, err := auth.GetAndValidateToken(r.Header, c.Secret)
	if err != nil {
		err = fmt.Errorf("error getting and validating token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(&req)
	if err != nil {
		err := fmt.Errorf("error decoding username request: %v", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	user, err := c.DbQueries.UpdateUserUsernameById(context.Background(), database.UpdateUserUsernameByIdParams{
		ID:        userId,
		Username:  req.Username,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		err = fmt.Errorf("error updating username to '%v' for '%v': %w", req.Username, userId, err)
		ResponseError(w, http.StatusConflict, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusAccepted, respUser{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (c *Config) HandlerUpdateUserEmailById(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email string `json:"email"`
	}
	userId, err := auth.GetAndValidateToken(r.Header, c.Secret)
	if err != nil {
		err = fmt.Errorf("error getting and validating token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(&req)
	if err != nil {
		err := fmt.Errorf("error decoding email request: %v", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	user, err := c.DbQueries.UpdateUserEmailById(context.Background(), database.UpdateUserEmailByIdParams{
		ID:        userId,
		Email:     req.Email,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		err = fmt.Errorf("error updating email to '%v' for '%v': %w", req.Email, userId, err)
		ResponseError(w, http.StatusConflict, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusAccepted, respUser{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (c *Config) HandlerUpdateUserPasswordById(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password string `json:"password"`
	}
	userId, err := auth.GetAndValidateToken(r.Header, c.Secret)
	if err != nil {
		err = fmt.Errorf("error getting and validating token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(&req)
	if err != nil {
		err := fmt.Errorf("error decoding password request: %v", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		err = fmt.Errorf("error hashing password for user '%v': %w", userId, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	user, err := c.DbQueries.UpdateUserPasswordById(context.Background(), database.UpdateUserPasswordByIdParams{
		ID:             userId,
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		err = fmt.Errorf("error updating password for '%v': %w", userId, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusAccepted, respUser{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
