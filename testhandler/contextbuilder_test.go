package testhandler

import (
	"context"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type ContextTestSuite struct {
	suite.Suite
}

func TestContextTestSuite(t *testing.T) {
	suite.Run(t, new(ContextTestSuite))
}

func (c *ContextTestSuite) TestCtx() {

	c.Run("a test contexts should be equal to itself", func() {
		// Run
		testCtx := Ctx(true, false)

		// Assert
		c.NotEmpty(testCtx.Value(echo.HeaderXRequestID))
		c.Equal(testCtx, testCtx)
	})

	c.Run("two test contexts should be unequal to each other", func() {
		// Run
		ctx1 := Ctx(false, false)
		ctx2 := Ctx(false, false)

		// Assert
		c.Empty(ctx1.Value(echo.HeaderXRequestID))
		c.Empty(ctx2.Value(echo.HeaderXRequestID))
		c.NotEqual(ctx1, ctx2)
	})

	c.Run("a test context should be unequal to a context.Background()", func() {
		// Run
		ctxBackground := context.Background()
		testCtx := Ctx(false, false)

		// Assert
		c.Empty(testCtx.Value(echo.HeaderXRequestID))
		c.NotEqual(ctxBackground, testCtx)
	})

	c.Run("a test context should be canceled", func() {
		// Run
		testCtx := Ctx(false, true)

		// Assert
		c.Empty(testCtx.Value(echo.HeaderXRequestID))
		c.Error(testCtx.Err())
	})
}
