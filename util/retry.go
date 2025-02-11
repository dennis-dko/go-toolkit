package util

import (
	"time"
)

func IncRetryDelay(attempt int, baseDelay time.Duration) time.Duration {
	if attempt <= 1 {
		return baseDelay
	}
	a, b := baseDelay, baseDelay
	for i := 2; i <= attempt; i++ {
		a, b = b, a+b
	}
	return b
}
