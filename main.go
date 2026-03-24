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
	mux.HandleFunc("POST api/users/login", config.HandlerLoginUser)
	mux.HandleFunc("POST api/users/register", config.HandlerRegisterUser)

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
