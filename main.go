package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/gabeamv/overloaded-api/api"
	"github.com/gabeamv/overloaded-api/internal/database"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env file: %v\n", err)
		os.Exit(1)
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening connection to database: %v\n", err)
		os.Exit(1)
	}
	secret := os.Getenv("SECRET")
	queries := database.New(db)
	config := api.Config{DbQueries: queries, Secret: secret}

	mux := http.NewServeMux()

	mux.HandleFunc("DELETE /admin/reset", config.HandlerDevReset)

	mux.HandleFunc("POST /api/users/login", config.HandlerLoginUser)
	mux.HandleFunc("POST /api/users/register", config.HandlerRegisterUser)

	mux.HandleFunc("POST /api/exercises", config.HandlerAddExercise)
	mux.HandleFunc("GET /api/exercises", config.HandlerGetAllCustomExercisesUserId)
	mux.HandleFunc("GET /api/exercises/{exercise_id}", config.HandlerGetExerciseById)
	mux.HandleFunc("PUT /api/exercises/{exercise_id}", config.HandlerUpdateExerciseNameById)
	mux.HandleFunc("DELETE /api/exercises/{exercise_id}", config.HandlerDeleteExerciseById)

	mux.HandleFunc("POST /api/workouts", config.HandlerAddWorkout)
	//mux.HandleFunc("PUT /api/workouts/{workout_id}", config.HandlerUpdateWorkoutById)
	mux.HandleFunc("GET /api/workouts", config.HandlerGetAllWorkouts)
	mux.HandleFunc("GET /api/workouts/{workout_id}", config.HandlerGetWorkoutById)
	mux.HandleFunc("DELETE /api/workouts/{workout_id}", config.HandlerDeleteWorkoutById)

	mux.HandleFunc("POST /api/workouts/{workout_id}/sets", config.HandlerAddSetByWorkoutId)
	//mux.HandleFunc("GET /api/workouts/{workout_id}/sets", config.HandlerGetSetsByWorkoutId)
	//mux.HandleFunc("GET /api/sets/{set_id}", config.HandlerGetSetById)
	//mux.HandleFunc("PUT /api/sets/{set_id}", config.HandlerUpdateSetById)
	//mux.HandleFunc("DELETE /api/sets/{set_id}", config.HandlerDeleteSetById)

	mux.HandleFunc("GET /api/progression_rules", config.HandlerGetAllProgressionRules)
	mux.HandleFunc("POST /api/progression_rules", config.HandlerDevApplyProgressionRules)
	mux.HandleFunc("DELETE /api/progression_rules", config.HandlerDevResetProgressionRules)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	err = server.ListenAndServe()
	if err != nil {
		err = fmt.Errorf("Server has closed: %w", err)
		log.Println(err)
	}
}
