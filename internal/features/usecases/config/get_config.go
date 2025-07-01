package config

import (
	"context"
	"strings"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
)

type getConfigUseCase struct{}

func NewGetConfigUseCase() GetConfigUseCase {
	return &getConfigUseCase{}
}

func (uc *getConfigUseCase) Execute(ctx context.Context, req GetConfigRequest) (*GetConfigResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, NewValidationError("invalid get config request", "", err)
	}

	// Read current configuration from file
	cfg, err := config.ReadConfigFile()
	if err != nil {
		return nil, NewFileSystemError("failed to read config file", err)
	}

	// Prepare response
	response := &GetConfigResponse{
		Config:   cfg,
		Values:   make(map[string]string),
		IsEmpty:  uc.isConfigEmpty(cfg),
		Warnings: []string{},
	}

	// If specific keys were requested, extract only those values
	if len(req.Keys) > 0 {
		for _, key := range req.Keys {
			value, exists := uc.getConfigValue(cfg, key)
			if exists {
				response.Values[key] = value
			} else {
				response.Warnings = append(response.Warnings, "Key '"+key+"' not found in configuration")
			}
		}
	} else {
		// Return all configuration values
		response.Values = uc.getAllConfigValues(cfg)
	}

	// Add configuration warnings
	warnings := uc.generateConfigWarnings(cfg)
	response.Warnings = append(response.Warnings, warnings...)

	return response, nil
}

func (uc *getConfigUseCase) getConfigValue(cfg config.AppConfig, key string) (string, bool) {
	// Normalize key to lowercase for comparison
	normalizedKey := strings.ToLower(key)

	switch normalizedKey {
	case "access_token", "accesstoken":
		return cfg.AccessToken, cfg.AccessToken != ""
	case "backend_url", "backendurl":
		return cfg.BackendURL, cfg.BackendURL != ""
	default:
		return "", false
	}
}

func (uc *getConfigUseCase) getAllConfigValues(cfg config.AppConfig) map[string]string {
	values := make(map[string]string)

	if cfg.AccessToken != "" {
		values["access_token"] = cfg.AccessToken
	}

	if cfg.BackendURL != "" {
		values["backend_url"] = cfg.BackendURL
	}

	return values
}

func (uc *getConfigUseCase) isConfigEmpty(cfg config.AppConfig) bool {
	return cfg.AccessToken == "" && cfg.BackendURL == ""
}

func (uc *getConfigUseCase) generateConfigWarnings(cfg config.AppConfig) []string {
	var warnings []string

	// Check for missing required configuration
	if cfg.AccessToken == "" {
		warnings = append(warnings, "Access token is not set. Use 'envsync config set access_token=<token>' to set it.")
	}

	if cfg.BackendURL == "" {
		warnings = append(warnings, "Backend URL is not set. Use 'envsync config set backend_url=<url>' to set it.")
	}

	// Check for insecure configurations
	if cfg.BackendURL != "" && strings.HasPrefix(cfg.BackendURL, "http://") {
		warnings = append(warnings, "Backend URL uses insecure HTTP protocol. Consider using HTTPS for better security.")
	}

	// Check for potentially invalid access token
	if cfg.AccessToken != "" && len(cfg.AccessToken) < 10 {
		warnings = append(warnings, "Access token appears to be too short and may be invalid.")
	}

	return warnings
}

func (uc *getConfigUseCase) maskSensitiveValue(key, value string) string {
	// Mask sensitive values for security
	normalizedKey := strings.ToLower(key)

	switch normalizedKey {
	case "access_token", "accesstoken":
		return uc.maskToken(value)
	default:
		return value
	}
}

func (uc *getConfigUseCase) maskToken(token string) string {
	if len(token) <= 8 {
		return strings.Repeat("*", len(token))
	}

	// Show first 4 and last 4 characters
	prefix := token[:4]
	suffix := token[len(token)-4:]
	middle := strings.Repeat("*", len(token)-8)

	return prefix + middle + suffix
}

func (uc *getConfigUseCase) validateConfigIntegrity(cfg config.AppConfig) error {
	// Perform basic integrity checks on the configuration

	// Check if both access token and backend URL are set for full functionality
	if cfg.AccessToken != "" && cfg.BackendURL == "" {
		return NewValidationError("access token is set but backend URL is missing", "backend_url", nil)
	}

	if cfg.BackendURL != "" && cfg.AccessToken == "" {
		return NewValidationError("backend URL is set but access token is missing", "access_token", nil)
	}

	return nil
}
