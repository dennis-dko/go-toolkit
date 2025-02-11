package envhandler

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"

	"github.com/joho/godotenv"
)

var envFiles = []string{
	".env.secrets.local",
	".env.local",
}

// Load loads the local files if exist and parse all environment variables into the given config struct
func Load(config interface{}) ([]string, error) {
	var loadedFiles []string

	// Load env file
	for _, envFile := range envFiles {
		err := godotenv.Load(envFile)
		if err != nil {
			// Skip env files which do not exist
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("failed to load env file: %w", err)
		}

		loadedFiles = append(loadedFiles, envFile)
	}

	// Parse env variables
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration from environment: %w", err)
	}

	return loadedFiles, nil
}
