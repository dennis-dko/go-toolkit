package util

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TlsTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (tl *TlsTestSuite) SetupTest() {
	tl.ctx = context.Background()
}

func TestTlsTestSuite(t *testing.T) {
	suite.Run(t, new(TlsTestSuite))
}

func (tl *TlsTestSuite) TestTlsConfig() {
	tl.Run("happy path - insecure connection", func() {
		// Run
		tlsInsecure, err := TlsConfig(tl.ctx, false, "", "", "")

		// Assert
		tl.Nil(tlsInsecure)
		tl.NoError(err)
	})
	tl.Run("happy path - secure connection error case", func() {
		// Run
		tlsSecure, err := TlsConfig(tl.ctx, true, "client.crt", "client.key", "ca.crt")

		// Assert
		tl.Error(err)
		tl.Nil(tlsSecure)
	})
}
