package envhandler

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Config struct {
	TestID       int    `env:"TEST_ID,unset"`
	TestName     string `env:"TEST_NAME,unset"`
	TestPassword string `env:"TEST_PASSWORD,unset"`
}

type EnvManagerTestSuite struct {
	suite.Suite
	cfg *Config
}

func (e *EnvManagerTestSuite) SetupSubTest() {
	// Sub setup
	e.cfg = &Config{}
}

func TestEnvManagerTestSuite(t *testing.T) {
	suite.Run(t, new(EnvManagerTestSuite))
}

func (e *EnvManagerTestSuite) TestLoad() {
	e.Run("happy path - return no error by loading from environment", func() {
		// Init
		e.T().Setenv("TEST_ID", "1")
		e.T().Setenv("TEST_NAME", "Tester")
		e.T().Setenv("TEST_PASSWORD", "querty")

		// Run
		loadedFiles, err := Load(e.cfg)

		// Assert
		e.NoError(err)
		e.Nil(loadedFiles)
		e.Equal(1, e.cfg.TestID)
		e.Equal("Tester", e.cfg.TestName)
		e.Equal("querty", e.cfg.TestPassword)
	})

	e.Run("happy path - return no error by loading from env files", func() {
		// Init
		envFiles = []string{
			"testdata/.env.secrets.test",
			"testdata/.env.test",
		}

		// Run
		loadedFiles, err := Load(e.cfg)

		// Assert
		e.NoError(err)
		e.Equal(envFiles, loadedFiles)
		e.Equal(1, e.cfg.TestID)
		e.Equal("Tester", e.cfg.TestName)
		e.Equal("querty", e.cfg.TestPassword)
	})

	e.Run("happy path - return no error if env file not exist", func() {
		// Init
		envFiles = []string{
			"testdata/.env.test.load.not.exist",
		}

		// Run
		loadedFiles, err := Load(e.cfg)

		// Assert
		e.NoError(err)
		e.Nil(loadedFiles)
	})

	e.Run("failed path - should return an error if the content of the file is not env", func() {
		// Init
		envFiles = []string{
			"testdata/.env.test.load.fail",
		}

		// Run
		loadedFiles, err := Load(e.cfg)

		// Assert
		e.Error(err)
		e.Nil(loadedFiles)
		e.ErrorContains(err, "failed to load env file: unexpected character \"?\" in variable name near \"TEST_NAME ? Tester\\n\"")
	})

	e.Run("failed path - should return an error if the file cannot be parsed into the config struct", func() {
		// Init
		envFiles = []string{
			"testdata/.env.test.parse.fail",
		}

		// Run
		loadedFiles, err := Load(e.cfg)

		// Assert
		e.Error(err)
		e.Nil(loadedFiles)
		e.ErrorContains(err, "failed to parse configuration from environment: env: parse error on field \"TestID\" of type \"int\": strconv.ParseInt: parsing \"true\": invalid syntax")
	})
}
