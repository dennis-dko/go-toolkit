package util

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RetryTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (r *RetryTestSuite) SetupTest() {
	r.ctx = context.Background()
}

func TestRetryTestSuite(t *testing.T) {
	suite.Run(t, new(RetryTestSuite))
}

func (r *RetryTestSuite) TestIncRetryDelay() {
	r.Run("happy path - increase retry delay", func() {
		// Run
		var result time.Duration
		for i := 1; i <= 10; i++ {
			retryDelay := 2 * time.Second
			retryDelay = IncRetryDelay(i, retryDelay)
			result += retryDelay
		}

		// Assert
		r.Equal(7.7, result.Minutes())
	})
}
