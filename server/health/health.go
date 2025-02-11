package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
}

func NewHealthController() *Controller {
	return &Controller{}
}

func (co *Controller) HandleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
}
