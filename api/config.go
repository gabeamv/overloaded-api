package api

import (
	"github.com/gabeamv/overloaded-api/internal/database"
)

type Config struct {
	// TODO: write attributes for database queries, secrets/keys later on from .env file
	DbQueries *database.Queries
	Secret    string
}
