package repository

import (
	"os"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
	"resty.dev/v3"
)

// createHTTPClient initializes and returns a new HTTP client with proper authentication
// and configuration for API requests.
func createHTTPClient() *resty.Client {
	// Initialize variables for authentication token and application configuration
	var authToken string
	var cfg config.AppConfig
	var cliCmd string

	// Check if API key is provided as an environment variable
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok || apiKey == "" {
		// If API_KEY environment variable is not set or empty,
		// load configuration from default config and use the access token
		cfg = config.New()
		authToken = cfg.AccessToken
	} else {
		// Otherwise use the API key from environment variable
		authToken = apiKey
	}

	// get the args passed to the CLI
	if len(os.Args) > 1 {
		cliCmd = os.Args[1]
	}

	// Create and configure a new REST client
	client := resty.New().
		SetDisableWarn(true).
		SetBaseURL(cfg.BackendURL).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-CLI-CMD", cliCmd).
		SetAuthToken(authToken)

	return client
}
