package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ErrorResponse struct {
	ErrorDetail struct {
		Code       int               `json:"code"`
		Message    map[string]string `json:"message"`
		Details    string            `json:"details"`
		Timestamp  time.Time         `json:"timestamp"`
		Path       string            `json:"path"`
		Suggestion string            `json:"suggestion"`
	} `json:"error"`
}

type ClientReporter interface {
	error
	Message() map[string]string
	StatusCode() int
	GetErrorDetais() string
	GetSuguestions() string
}

func responseWithErrors(w http.ResponseWriter, r http.Request, err error) {
	if cr, ok := err.(ClientReporter); ok {
		status := cr.StatusCode()
		if status >= http.StatusInternalServerError {
			responseWithError(w, status, "Something has gone wrong")
		}

		errRep := ErrorResponse{
			ErrorDetail: struct {
				Code       int               `json:"code"`
				Message    map[string]string `json:"message"`
				Details    string            `json:"details"`
				Timestamp  time.Time         `json:"timestamp"`
				Path       string            `json:"path"`
				Suggestion string            `json:"suggestion"`
			}{
				Code:       cr.StatusCode(),
				Message:    cr.Message(),
				Details:    cr.GetErrorDetais(),
				Timestamp:  time.Now().UTC(),
				Path:       r.URL.RequestURI(),
				Suggestion: cr.GetSuguestions(),
			},
		}

		respondWithJson(w, errRep.ErrorDetail.Code, errRep)
		return
	}

	responseWithError(w, 500, "Something has gone wrong")
}
func responseWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error:", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
