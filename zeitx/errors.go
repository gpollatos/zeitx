package zeitx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type apiError struct {
	StatusCode int       `json:"statusCode"`
	StatusMsg  string    `json:"statusMsg"`
	Details    string    `json:"details"`
	Ts         time.Time `json:"ts"`
}

func (e apiError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		var GenericError = apiError{
			StatusCode: http.StatusInternalServerError,
			StatusMsg:  http.StatusText(http.StatusInternalServerError),
			Details:    err.Error(),
			Ts:         time.Now(),
		}
		b, _ = json.Marshal(GenericError)
	}
	return string(b)
}

// APIError writes a json formatted error to the ResponseWriter
func APIError(w http.ResponseWriter, status int, details string) {
	w.Header().Set(ContentType, ApplicationJSONCharsetUTF8)
	w.Header().Set(XContentTypeOptions, "nosniff")
	err := apiError{
		StatusCode: status,
		StatusMsg:  http.StatusText(status),
		Details:    details,
		Ts:         time.Now(),
	}
	w.WriteHeader(err.StatusCode)
	fmt.Fprintln(w, err)
}

// BadRequest is a helper method for sending a 400
func BadRequest(w http.ResponseWriter, details string) {
	APIError(w, http.StatusBadRequest, details)
}

// Unauthorized is a helper method for sending a 401
func Unauthorized(w http.ResponseWriter, details string) {
	APIError(w, http.StatusUnauthorized, details)
}

// NotFound is a helper method for sending a 404
func NotFound(w http.ResponseWriter, details string) {
	APIError(w, http.StatusNotFound, details)
}

// InternalServerError is a helper method for sending a 500
func InternalServerError(w http.ResponseWriter, details string) {
	APIError(w, http.StatusInternalServerError, details)
}
