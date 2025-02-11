package util

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UUIDSuite struct {
	suite.Suite
}

func TestUUIDSuite(t *testing.T) {
	suite.Run(t, new(UUIDSuite))
}

func (u *UUIDSuite) TestSetUUID() {

	u.Run("happy path - set uuid", func() {
		// Run
		uuid := SetUUID()

		// Assert
		u.NotEmpty(uuid)
	})
}
