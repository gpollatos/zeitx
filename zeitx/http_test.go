package zeitx_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gpollatos/zeitx-all/zeitx"

	"github.com/stretchr/testify/assert"
)

func TestOkJSON(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := NewInMemResponse()
	payload := &response{Key: "test"}
	zeitx.OkJSON(w, r, payload)
	assertJSONResponse(t, w, 200, *payload)
}

func TestJSON(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := NewInMemResponse()
	payload := &response{Key: "test"}
	zeitx.JSON(w, r, payload, 201)
	assertJSONResponse(t, w, 201, *payload)
}

type response struct {
	Key string
}

func assertJSONResponse(t *testing.T, w *inMemResponse, status int, body response) {
	assert.Equal(t, status, w.statusCode)
	assert.Contains(t, w.Header(), zeitx.ContentType)
	assert.Equal(t, []string{zeitx.ApplicationJSONCharsetUTF8}, w.Header()[zeitx.ContentType])

	r := response{}
	err := json.Unmarshal(w.body, &r)
	if err != nil {
		t.Errorf("the response should be a valid json with the expected structure: %s", err)
	}

	assert.Equal(t, status, w.statusCode)
	assert.Equal(t, body, r)
}
