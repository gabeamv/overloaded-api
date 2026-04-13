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

type respWorkoutSet struct {
	Id            uuid.UUID `json:"id"`
	WorkoutId     uuid.UUID `json:"workout_id"`
	ExerciseId    uuid.UUID `json:"exercise_id"`
	ProgressTrack uuid.UUID `json:"progress_track"`
	WeightInLbs   float64   `json:"weight_in_lbs"`
	Reps          float64   `json:"reps"`
	TimeInSeconds float64   `json:"time_in_seconds"`
}

func (c *Config) HandlerAddSetByWorkoutId(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ExerciseId    uuid.UUID `json:"exercise_id"`
		ProgressTrack uuid.UUID `json:"progress_track"`
		WeightInLbs   float64   `json:"weight_in_lbs"`
		Reps          float64   `json:"reps"`
		TimeInSeconds float64   `json:"time_in_seconds"`
	}
	userId, err := auth.GetAndValidateToken(r.Header, c.Secret)
	if err != nil {
		err = fmt.Errorf("error getting and validating token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	workoutId, err := uuid.Parse(r.PathValue("workout_id"))
	if err != nil {
		err = fmt.Errorf("error parsing path value 'workout_id' into type UUID: %w", err)
		ResponseError(w, http.StatusNotFound, err.Error(), err)
		return
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(&req)
	if err != nil {
		err := fmt.Errorf("error decoding workout set request: %v", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	workout, err := c.DbQueries.GetWorkoutById(context.Background(), workoutId)
	if err != nil {
		err = fmt.Errorf("error getting workout '%v': %w", workoutId, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if userId != workout.UserID {
		err = fmt.Errorf("error, unauthorized request for user '%v' to add set to workout belonging to user '%v'", userId, workout.UserID)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	workoutSet, err := c.DbQueries.CreateWorkoutSetByWorkoutId(context.Background(), database.CreateWorkoutSetByWorkoutIdParams{
		WorkoutID:     workoutId,
		ExerciseID:    req.ExerciseId,
		ProgressTrack: req.ProgressTrack,
		WeightInLbs:   req.WeightInLbs,
		Reps:          req.Reps,
		TimeInSeconds: req.TimeInSeconds,
	})
	if err != nil {
		err = fmt.Errorf("error creating workout set for workout '%v': %w", workoutId, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusCreated, respWorkoutSet{
		Id:            workoutSet.ID,
		WorkoutId:     workoutSet.WorkoutID,
		ExerciseId:    workoutSet.ExerciseID,
		ProgressTrack: workoutSet.ProgressTrack,
		WeightInLbs:   workoutSet.WeightInLbs,
		Reps:          workoutSet.Reps,
		TimeInSeconds: workoutSet.TimeInSeconds,
	})
}

func (c *Config) HandlerGetSetsByWorkoutId(w http.ResponseWriter, r *http.Request) {

}
