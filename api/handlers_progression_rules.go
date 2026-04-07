package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gabeamv/overloaded-api/internal/auth"
	"github.com/gabeamv/overloaded-api/internal/database"
	"github.com/google/uuid"
)

type progressionResp struct {
	Id          uuid.UUID `json:"id"`
	Label       string    `json:"label"`
	Rule        string    `json:"rule"`
	Description string    `json:"description"`
}

func (c *Config) HandlerDevApplyProgressionRules(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		msg := "non developer api call"
		err := fmt.Errorf("%v", msg)
		ResponseError(w, http.StatusInternalServerError, msg, err)
		return
	}
	err := c.DbQueries.CreateProgressionRule(context.Background(), database.CreateProgressionRuleParams{
		Label:       "p",
		Rule:        "progress",
		Description: "Progress and go up in weight next session.",
	})
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	err = c.DbQueries.CreateProgressionRule(context.Background(), database.CreateProgressionRuleParams{
		Label:       "s",
		Rule:        "stay",
		Description: "Stay at the current weight you are at for next session.",
	})
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	err = c.DbQueries.CreateProgressionRule(context.Background(), database.CreateProgressionRuleParams{
		Label: "t",
		Rule:  "tag on form",
		Description: "Did not complete working set in the higher end of required rep range with perfect form." +
			" Stay at the current weight you are at for next session.",
	})
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusNoContent, struct{}{})
}

func (c *Config) HandlerDevResetProgressionRules(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		msg := "non developer api call"
		err := fmt.Errorf("%v", msg)
		ResponseError(w, http.StatusInternalServerError, msg, err)
		return
	}
	err := c.DbQueries.DeleteAllProgressionRules(context.Background())
	if err != nil {
		err = fmt.Errorf("error deleting all progression rules: %w", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
}

func (c *Config) HandlerGetAllProgressionRules(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err = fmt.Errorf("error getting token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	_, err = auth.ValidateJWT(token, c.Secret)
	if err != nil {
		err = fmt.Errorf("unauthorized request to get all progression rules: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	rules, err := c.DbQueries.GetAllProgressionRules(context.Background())
	if err != nil {
		err = fmt.Errorf("error getting all progression rules: %w", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	var resp []progressionResp
	for _, rule := range rules {
		resp = append(resp, progressionResp{
			Id:          rule.ID,
			Label:       rule.Label,
			Rule:        rule.Rule,
			Description: rule.Description,
		})
	}
	ResponseJSON(w, http.StatusAccepted, resp)
}
