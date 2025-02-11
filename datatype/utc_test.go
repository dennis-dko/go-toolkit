package datatype

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtcSuite struct {
	suite.Suite
}

func TestUtcSuite(t *testing.T) {
	suite.Run(t, new(UtcSuite))
}

func (u *UtcSuite) TestIsUtcTime() {

	u.Run("happy path - datetime with utc timezone", func() {
		// Run
		utc := IsUtcTime("2024-06-13T12:00:00Z")

		// Assert
		u.True(utc)
	})
	u.Run("happy path - datetime without utc timezone", func() {
		// Run
		utc := IsUtcTime("2024-06-13T12:00:00")

		// Assert
		u.False(utc)
	})
}

func (u *UtcSuite) TestUtcSuffixSet() {

	u.Run("happy path - datetime set utc suffix", func() {
		// Run
		utc := SetAsUtc("2024-06-13T12:00:00")

		// Assert
		u.Equal("2024-06-13T12:00:00Z", utc)
	})
}
