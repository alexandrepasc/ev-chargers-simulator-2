package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	ep  string = "/health"
	iep string = "/invalid"
)

func TestApiInstance(t *testing.T) {
	r, w := beforeEach()

	req, _ := http.NewRequest("GET", ep, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNotFound(t *testing.T) {
	r, w := beforeEach()

	req, _ := http.NewRequest("GET", iep, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	assert.Equal(t, "{\"message\":\"Route not found.\"}", w.Body.String())
}

func TestNotAllowed(t *testing.T) {
	r, w := beforeEach()

	req, _ := http.NewRequest("POST", ep, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	assert.Equal(t, "{\"message\":\"Method not allowed.\"}", w.Body.String())
}

func beforeEach() (r *gin.Engine, w *httptest.ResponseRecorder) {
	r = api.New()

	w = httptest.NewRecorder()

	return r, w
}
