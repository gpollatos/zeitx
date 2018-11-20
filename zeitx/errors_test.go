package zeitx_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gpollatos/zeitx/zeitx"

	"github.com/stretchr/testify/assert"
)

func TestAPIErrorWithA200(t *testing.T) {
	w := NewInMemResponse()
	zeitx.APIError(w, 200, "OK")
	assertErrorResponse(t, w, 200, "OK")
}

func TestBadRequest(t *testing.T) {
	w := NewInMemResponse()
	zeitx.BadRequest(w, "some message")
	assertErrorResponse(t, w, 400, "some message")
}

func TestUnauthorized(t *testing.T) {
	w := NewInMemResponse()
	zeitx.Unauthorized(w, "some message")
	assertErrorResponse(t, w, 401, "some message")
}

func TestNotFound(t *testing.T) {
	w := NewInMemResponse()
	zeitx.NotFound(w, "some message")
	assertErrorResponse(t, w, 404, "some message")
}

func TestInternalServerError(t *testing.T) {
	w := NewInMemResponse()
	zeitx.InternalServerError(w, "some message")
	assertErrorResponse(t, w, 500, "some message")
}

func assertErrorResponse(t *testing.T, w *inMemResponse, status int, details string) {
	assert.Equal(t, status, w.statusCode)
	assert.Contains(t, w.Header(), zeitx.ContentType)
	assert.Contains(t, w.Header(), zeitx.XContentTypeOptions)
	assert.Equal(t, []string{zeitx.ApplicationJSONCharsetUTF8}, w.Header()[zeitx.ContentType])
	assert.Equal(t, []string{"nosniff"}, w.Header()[zeitx.XContentTypeOptions])

	e := errorResponse{}
	err := json.Unmarshal(w.body, &e)
	if err != nil {
		t.Errorf("the response should be a valid json with the expected structure: %s", err)
	}

	assert.Equal(t, status, e.StatusCode)
	assert.Equal(t, http.StatusText(status), e.StatusMsg)
	assert.Equal(t, details, e.Details)
}

type errorResponse struct {
	StatusCode int
	StatusMsg  string
	Details    string
	Ts         time.Time
}

type inMemResponse struct {
	headers    http.Header
	body       []byte
	statusCode int
}

func NewInMemResponse() *inMemResponse {
	return &inMemResponse{
		headers: make(http.Header),
	}
}

func (w *inMemResponse) Header() http.Header {
	return w.headers
}

func (w *inMemResponse) Write(body []byte) (int, error) {
	w.body = make([]byte, len(body))
	return copy(w.body, body), nil
}

func (w *inMemResponse) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
