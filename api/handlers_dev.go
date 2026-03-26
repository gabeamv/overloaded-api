package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

func (c *Config) HandlerDevReset(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		msg := "non developer api call"
		err := fmt.Errorf("%v", msg)
		ResponseError(w, http.StatusInternalServerError, msg, err)
		return
	}
	err := c.DbQueries.DeleteAllUsers(context.Background())
	if err != nil {
		msg := "error deleting all users"
		err = fmt.Errorf("%v: %w", msg, err)
		ResponseError(w, http.StatusInternalServerError, msg, err)
		return
	}
	resp := "All users deleted."
	ResponseJSON(w, http.StatusOK, resp)
}
