package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/health"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const ep string = "/health"

func TestGetHealth(t *testing.T) {
	r := gin.Default()

	health.Health(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", ep, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, "{\"message\":\"OK\"}", w.Body.String())
}
