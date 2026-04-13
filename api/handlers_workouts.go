package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabeamv/overloaded-api/internal/auth"
	"github.com/gabeamv/overloaded-api/internal/database"
	"github.com/google/uuid"
)

type workoutResp struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	Volume    int32     `json:"volume"`
	Prs       int32     `json:"prs"`
}

func (c *Config) HandlerAddWorkout(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserId    uuid.UUID `json:"user_id"`
		StartedAt time.Time `json:"started_at"`
		EndedAt   time.Time `json:"ended_at"`
		Volume    int32     `json:"volume"`
		Prs       int32     `json:"prs"`
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
		err := fmt.Errorf("error decoding workout request: %v", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if userId != req.UserId {
		err = fmt.Errorf("error, unauthorized request for user '%v' to add workout for user '%v'.", userId, req.UserId)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	workout, err := c.DbQueries.CreateWorkout(context.Background(), database.CreateWorkoutParams{
		UserID:    req.UserId,
		StartedAt: req.StartedAt,
		EndedAt:   req.EndedAt,
		Volume:    req.Volume,
		Prs:       req.Prs,
	})
	if err != nil {
		err = fmt.Errorf("error creating workout for user '%v'.", userId)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	resp := workoutResp{
		Id:        workout.ID,
		UserId:    workout.UserID,
		StartedAt: workout.StartedAt,
		EndedAt:   workout.EndedAt,
		Volume:    workout.Volume,
		Prs:       workout.Prs,
	}
	ResponseJSON(w, http.StatusOK, resp)
}

func (c *Config) HandlerGetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetAndValidateToken(r.Header, c.Secret)
	if err != nil {
		err = fmt.Errorf("error getting and validating token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	workouts, err := c.DbQueries.GetAllWorkoutsByUserId(context.Background(), userId)
	if err != nil {
		err = fmt.Errorf("error getting all workouts for user '%v': %w", userId, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	var resp []workoutResp
	for _, workout := range workouts {
		resp = append(resp, workoutResp{
			Id:        workout.ID,
			UserId:    workout.UserID,
			StartedAt: workout.StartedAt,
			EndedAt:   workout.EndedAt,
			Volume:    workout.Volume,
			Prs:       workout.Prs,
		})
	}
	ResponseJSON(w, http.StatusAccepted, resp)
}

func (c *Config) HandlerGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetAndValidateToken(r.Header, c.Secret)
	if err != nil {
		err = fmt.Errorf("error getting and validating token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	id, err := uuid.Parse(r.PathValue("workout_id"))
	if err != nil {
		err = fmt.Errorf("error parsing path value 'workout_id' into type UUID: %w", err)
		ResponseError(w, http.StatusNotFound, err.Error(), err)
		return
	}
	workout, err := c.DbQueries.GetWorkoutById(context.Background(), id)
	if err != nil {
		err = fmt.Errorf("error getting workout for user '%v': %w", userId, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusAccepted, workoutResp{
		Id:        workout.ID,
		UserId:    workout.UserID,
		StartedAt: workout.StartedAt,
		EndedAt:   workout.EndedAt,
		Volume:    workout.Volume,
		Prs:       workout.Prs,
	})
}

func (c *Config) HandlerDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetAndValidateToken(r.Header, c.Secret)
	if err != nil {
		err = fmt.Errorf("error getting and validating token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	id, err := uuid.Parse(r.PathValue("workout_id"))
	if err != nil {
		err = fmt.Errorf("error parsing path value 'workout_id' into type UUID: %w", err)
		ResponseError(w, http.StatusNotFound, err.Error(), err)
		return
	}
	err = c.DbQueries.DeleteWorkoutById(context.Background(), id)
	if err != nil {
		err = fmt.Errorf("error deleting workout '%v' for user '%v': %w", id, userId, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusNoContent, struct{}{})
}
