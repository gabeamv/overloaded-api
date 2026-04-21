package api

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/gabeamv/overloaded-api/internal/auth"
	"github.com/gabeamv/overloaded-api/internal/database"
	"github.com/google/uuid"
)

// Calculates estimated One Rep Max of a specific exercise using the Epley equation.
func (c *Config) HandlerGetPrOneRepMax(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		OneRepMax float64 `json:"one_rep_max"`
	}
	userId, err := auth.GetAndValidateToken(r.Header, c.Secret)
	if err != nil {
		err = fmt.Errorf("error getting and validating token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	exerciseId, err := uuid.Parse(r.PathValue("exercise_id"))
	if err != nil {
		err = fmt.Errorf("error parsing path value 'exercise_id' into type UUID: %w", err)
		ResponseError(w, http.StatusNotFound, err.Error(), err)
		return
	}
	oneRepMax, err := c.DbQueries.GetPrOneRepMaxByExerciseIdUserId(context.Background(), database.GetPrOneRepMaxByExerciseIdUserIdParams{
		ExerciseID: exerciseId,
		UserID:     userId,
	})
	if err != nil {
		err = fmt.Errorf("error getting one rep max for user '%v' for exercise '%v': %w", userId, exerciseId, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusAccepted, resp{OneRepMax: oneRepMax})
}

// Calculates most volume done for a set of a specific exercise.
func (c *Config) HandlerGetPrVolume(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		VolumeMax float64 `json:"volume_max"`
	}
	userId, err := auth.GetAndValidateToken(r.Header, c.Secret)
	if err != nil {
		err = fmt.Errorf("error getting and validating token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	exerciseId, err := uuid.Parse(r.PathValue("exercise_id"))
	if err != nil {
		err = fmt.Errorf("error parsing path value 'exercise_id' into type UUID: %w", err)
		ResponseError(w, http.StatusNotFound, err.Error(), err)
		return
	}
	volumeMax, err := c.DbQueries.GetPrVolumeByExerciseIdUserId(context.Background(), database.GetPrVolumeByExerciseIdUserIdParams{
		ExerciseID: exerciseId,
		UserID:     userId,
	})
	if err != nil {
		err = fmt.Errorf("error getting volume max for user '%v' for exercise '%v': %w", userId, exerciseId, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusAccepted, resp{VolumeMax: volumeMax})
}

func roundToTwo(x float64) float64 {
	return math.Round(x*100) / 100
}
