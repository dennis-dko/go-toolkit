package testhandler

import (
	"context"

	"github.com/dennis-dko/go-toolkit/util"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/rand"
)

const contextKey = "TestContextUniqueValue"

// Ctx returns a context to be used in tests.
//
// The difference to context.Background() is as follows:
// Two contexts produced with context.Background() will be equal, two context produced with this method will not be
// equal. Also, a context produced with this method will always be unequal to one produced with context.Background().
//
// Problem this solves:
// A test should be able to detect if a context has been manipulated when it calls mock functions. Most times the
// context should just be propagated. If the code to be tested is called with a context.Background() and also the call
// to the mock function is accidentally done with a new context.Background(), the test will not be able to detect this.
// If using this function instead to generate the context to call the code to be tested, the test can appropriately
// fail when the contexts are checked for equality.
func Ctx(useReqID bool, isCancel bool) context.Context {
	testCtx := context.WithValue(context.Background(), contextKey, rand.Int63())
	if useReqID {
		testCtx = context.WithValue(testCtx, echo.HeaderXRequestID, util.SetUUID())
	}
	if isCancel {
		cancelCtx, cancelFunc := context.WithCancel(testCtx)
		cancelFunc()
		testCtx = cancelCtx
	}
	return testCtx
}
