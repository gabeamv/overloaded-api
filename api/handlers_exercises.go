package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gabeamv/overloaded-api/internal/auth"
	"github.com/gabeamv/overloaded-api/internal/database"
	"github.com/google/uuid"
)

type exerciseResp struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	IsCustom bool      `json:"is_custom"`
	UserId   uuid.UUID `json:"user_id"`
}

func (c *Config) HandlerAddExercise(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string `json:"name"`
		IsCustom bool   `json:"is_custom"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err = fmt.Errorf("error getting token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	userId, err := auth.ValidateJWT(token, c.Secret)
	if err != nil {
		err = fmt.Errorf("unauthorized request to add exercise: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	var bodyExercise request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(&bodyExercise)
	if err != nil {
		err := fmt.Errorf("error decoding exercise: %v", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	exercise, err := c.DbQueries.CreateExercise(context.Background(), database.CreateExerciseParams{Name: bodyExercise.Name,
		IsCustom: bodyExercise.IsCustom, UserID: uuid.NullUUID{UUID: userId, Valid: true}})
	if err != nil {
		err = fmt.Errorf("error adding exercise: %w", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	resp := exerciseResp{Id: exercise.ID, Name: exercise.Name, IsCustom: exercise.IsCustom, UserId: exercise.UserID.UUID}
	ResponseJSON(w, http.StatusOK, resp)
}

func (c *Config) HandlerGetAllCustomExercisesUserId(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err = fmt.Errorf("error getting token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	userId, err := auth.ValidateJWT(token, c.Secret)
	if err != nil {
		err = fmt.Errorf("unauthorized request to get users custom exercises: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	exercises, err := c.DbQueries.GetAllExercisesUserId(context.Background(), uuid.NullUUID{UUID: userId, Valid: true})
	if err != nil {
		err = fmt.Errorf("error getting all users' custom exercises: %w", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	resp := make([]exerciseResp, 0)
	for _, exercise := range exercises {
		resp = append(resp, exerciseResp{Id: exercise.ID, Name: exercise.Name, IsCustom: exercise.IsCustom, UserId: exercise.UserID.UUID})
	}
	ResponseJSON(w, http.StatusOK, resp)
}
