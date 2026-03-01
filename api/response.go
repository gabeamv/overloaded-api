package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		log.Printf("an error has occured: %v\n", err)
	}
	if statusCode > 499 {
		log.Printf("responding with error code %v: %v\n", statusCode, message)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	ResponseJSON(w, statusCode, errResponse{Error: message})
}

func ResponseJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling '%v' to bytes: %v\n", payload, err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}
