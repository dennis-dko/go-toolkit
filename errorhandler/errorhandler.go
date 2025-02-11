package errorhandler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HttpErrorHandler struct {
	statusCodes map[error]int
}

// New creates a new HttpErrorHandler
func New(errorStatusCodeMaps map[error]int) *HttpErrorHandler {
	return &HttpErrorHandler{
		statusCodes: errorStatusCodeMaps,
	}
}

// Handler handles the error and sends the response to the client
func (h *HttpErrorHandler) Handler(err error, c echo.Context) {
	var he *echo.HTTPError
	ok := errors.As(err, &he)
	if ok {
		if he.Internal != nil {
			var httpErr *echo.HTTPError
			if errors.As(he.Internal, &httpErr) {
				he = httpErr
			}
		}
	} else {
		he = &echo.HTTPError{
			Code:    h.getStatusCode(err),
			Message: unwrapRecursive(err).Error(),
		}
	}
	code := he.Code
	message := he.Message
	if _, ok := he.Message.(string); ok {
		message = map[string]interface{}{"message": err.Error()}
	}
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	}
}

func (h *HttpErrorHandler) getStatusCode(err error) int {
	for key, value := range h.statusCodes {
		if errors.Is(err, key) {
			return value
		}
	}
	return http.StatusInternalServerError
}

func unwrapRecursive(err error) error {
	originalErr := err
	for originalErr != nil {
		var internalErr = errors.Unwrap(originalErr)
		if internalErr == nil {
			break
		}
		originalErr = internalErr
	}
	return originalErr
}
