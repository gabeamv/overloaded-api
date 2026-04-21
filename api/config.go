package api

import (
	"database/sql"

	"github.com/gabeamv/overloaded-api/internal/database"
)

type Config struct {
	Db        *sql.DB
	DbQueries *database.Queries
	Secret    string
}
