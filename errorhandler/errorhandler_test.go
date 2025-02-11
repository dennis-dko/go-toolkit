package errorhandler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

var ErrTestFailed = errors.New("test error")

type ErrorhandlerTestSuite struct {
	suite.Suite
	new              *HttpErrorHandler
	context          echo.Context
	recorder         *httptest.ResponseRecorder
	request          *http.Request
	expectedResponse string
}

func (e *ErrorhandlerTestSuite) SetupTest() {
	// Setup
	rec := httptest.NewRecorder()
	e.recorder = rec
	req := httptest.NewRequest(http.MethodPost, "http://localhost", nil)
	e.request = req
	ec := echo.New()
	echoContext := ec.NewContext(req, rec)
	e.context = echoContext
	statusCodeMap := NewErrorStatusCodeMaps()
	statusCodeMap[ErrTestFailed] = http.StatusBadRequest
	e.new = New(statusCodeMap)
	e.expectedResponse = `{"message":"test error"}`
}

func TestErrorhandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorhandlerTestSuite))
}

func (e *ErrorhandlerTestSuite) TestErrorHandling() {

	e.Run("happy path - get correct status code and error", func() {
		// Run
		e.new.Handler(ErrTestFailed, e.context)

		// Assert
		e.Equal(http.StatusBadRequest, e.recorder.Code)
		e.Equal(e.expectedResponse+"\n", e.recorder.Body.String())
	})

	e.Run("happy path - get only correct status", func() {
		// Init
		e.request.Method = http.MethodHead

		// Run
		e.new.Handler(ErrTestFailed, e.context)

		// Assert
		e.Equal(http.StatusBadRequest, e.recorder.Code)
	})
}
