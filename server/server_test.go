package server

import (
	"os"
	"testing"
	"time"

	"github.com/dennis-dko/go-toolkit/testhandler"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	newServer *Server
}

func (s *ServerTestSuite) SetupTest() {
	// Setup
	serverConfig := &Config{
		Host:                    "127.0.0.1",
		Port:                    0,
		GracefulShutDownTimeout: time.Minute,
	}
	s.newServer = New(testhandler.Ctx(false, false), serverConfig, nil)
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) TestStart() {

	s.Run("happy path - server is starting", func() {
		// Init
		var signalChannel chan os.Signal
		signalNotify = func(c chan os.Signal, sig ...os.Signal) {
			signalChannel = c
			signalChannel <- os.Interrupt
		}

		// Run
		s.newServer.Start()

		// Assert that the test finishes (and is not stuck in Start())
	})
}
