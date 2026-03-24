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

func (c *Config) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Id        uuid.UUID `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
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
	resp := response{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, Username: req.Username}
	ResponseJSON(w, http.StatusCreated, resp)
}

func (c *Config) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {

}
