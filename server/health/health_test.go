package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type HealthTestSuite struct {
	suite.Suite
}

func TestHealthTestSuite(t *testing.T) {
	suite.Run(t, new(HealthTestSuite))
}

func (h *HealthTestSuite) TestHealth() {

	h.Run("a test contexts should be equal to itself", func() {
		// Init
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		healthController := NewHealthController()

		// Run
		_ = healthController.HandleHealth(c)
		responseBodyMap := make(map[string]interface{})
		err := json.Unmarshal(rec.Body.Bytes(), &responseBodyMap)
		statusField, fieldOk := responseBodyMap["status"]
		statusValue, valueOk := statusField.(string)

		// Assert
		h.NoError(err)
		h.True(fieldOk)
		h.True(valueOk)
		h.Equal(http.StatusOK, rec.Result().StatusCode)
		h.Equal("UP", statusValue)
	})
}
