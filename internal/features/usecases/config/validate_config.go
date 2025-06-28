package config

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
)

type validateConfigUseCase struct{}

func NewValidateConfigUseCase() ValidateConfigUseCase {
	return &validateConfigUseCase{}
}

func (uc *validateConfigUseCase) Execute(ctx context.Context, req ValidateConfigRequest) (*ValidateConfigResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, NewValidationError("invalid validate config request", "", err)
	}

	var cfg config.AppConfig
	var err error

	// Use provided config or read from file
	if req.Config != nil {
		cfg = *req.Config
	} else {
		cfg, err = config.ReadConfigFile()
		if err != nil {
			return &ValidateConfigResponse{
				IsValid: false,
				Issues: []ValidationIssue{
					{
						Key:        "config_file",
						Message:    "Failed to read configuration file",
						Severity:   ValidationSeverityError,
						Suggestion: "Run 'envsync config set' to create a new configuration",
					},
				},
			}, nil
		}
	}

	// Perform validation
	issues := uc.validateConfiguration(ctx, cfg)

	// Determine if configuration is valid (no error-level issues)
	isValid := true
	for _, issue := range issues {
		if issue.Severity == ValidationSeverityError {
			isValid = false
			break
		}
	}

	return &ValidateConfigResponse{
		IsValid: isValid,
		Issues:  issues,
	}, nil
}

func (uc *validateConfigUseCase) validateConfiguration(ctx context.Context, cfg config.AppConfig) []ValidationIssue {
	var issues []ValidationIssue

	// Validate access token
	tokenIssues := uc.validateAccessToken(cfg.AccessToken)
	issues = append(issues, tokenIssues...)

	// Validate backend URL
	urlIssues := uc.validateBackendURL(ctx, cfg.BackendURL)
	issues = append(issues, urlIssues...)

	// Validate configuration completeness
	completenessIssues := uc.validateConfigCompleteness(cfg)
	issues = append(issues, completenessIssues...)

	// Validate configuration consistency
	consistencyIssues := uc.validateConfigConsistency(cfg)
	issues = append(issues, consistencyIssues...)

	return issues
}

func (uc *validateConfigUseCase) validateAccessToken(token string) []ValidationIssue {
	var issues []ValidationIssue

	if token == "" {
		issues = append(issues, ValidationIssue{
			Key:        "access_token",
			Message:    MsgAccessTokenRequired,
			Severity:   ValidationSeverityError,
			Suggestion: SuggestRunLogin,
		})
		return issues
	}

	// Check token length
	if len(token) < 10 {
		issues = append(issues, ValidationIssue{
			Key:        "access_token",
			Message:    "Access token is too short (minimum 10 characters)",
			Severity:   ValidationSeverityError,
			Suggestion: "Verify your access token is complete and valid",
		})
	}

	// Check token format (basic validation)
	if !uc.isValidTokenFormat(token) {
		issues = append(issues, ValidationIssue{
			Key:        "access_token",
			Message:    MsgAccessTokenWeak,
			Severity:   ValidationSeverityWarning,
			Suggestion: "Ensure you're using a valid access token from EnvSync Cloud",
		})
	}

	// Check for common token issues
	if strings.Contains(token, " ") {
		issues = append(issues, ValidationIssue{
			Key:        "access_token",
			Message:    "Access token contains spaces",
			Severity:   ValidationSeverityError,
			Suggestion: "Remove any spaces from your access token",
		})
	}

	return issues
}

func (uc *validateConfigUseCase) validateBackendURL(ctx context.Context, url string) []ValidationIssue {
	var issues []ValidationIssue

	if url == "" {
		issues = append(issues, ValidationIssue{
			Key:        "backend_url",
			Message:    MsgBackendURLRequired,
			Severity:   ValidationSeverityError,
			Suggestion: SuggestSetBackendURL,
		})
		return issues
	}

	// Check URL format
	if !uc.isValidURL(url) {
		issues = append(issues, ValidationIssue{
			Key:        "backend_url",
			Message:    "Backend URL format is invalid",
			Severity:   ValidationSeverityError,
			Suggestion: "Use format: https://api.example.com",
		})
		return issues
	}

	// Check for insecure HTTP
	if strings.HasPrefix(url, "http://") {
		issues = append(issues, ValidationIssue{
			Key:        "backend_url",
			Message:    MsgBackendURLInsecure,
			Severity:   ValidationSeverityWarning,
			Suggestion: SuggestUseHTTPS,
		})
	}

	// Test connectivity (optional, with timeout)
	connectivityIssue := uc.testConnectivity(ctx, url)
	if connectivityIssue != nil {
		issues = append(issues, *connectivityIssue)
	}

	return issues
}

func (uc *validateConfigUseCase) validateConfigCompleteness(cfg config.AppConfig) []ValidationIssue {
	var issues []ValidationIssue

	// Check if configuration is completely empty
	if cfg.AccessToken == "" && cfg.BackendURL == "" {
		issues = append(issues, ValidationIssue{
			Key:        "config",
			Message:    "Configuration is empty",
			Severity:   ValidationSeverityError,
			Suggestion: SuggestRunLogin,
		})
		return issues
	}

	// Check for incomplete configuration
	if cfg.AccessToken == "" || cfg.BackendURL == "" {
		issues = append(issues, ValidationIssue{
			Key:        "config",
			Message:    MsgConfigIncomplete,
			Severity:   ValidationSeverityWarning,
			Suggestion: "Set both access_token and backend_url for full functionality",
		})
	}

	return issues
}

func (uc *validateConfigUseCase) validateConfigConsistency(cfg config.AppConfig) []ValidationIssue {
	var issues []ValidationIssue

	// Check for mismatched configuration
	// For example, if using a development token with production URL
	if cfg.AccessToken != "" && cfg.BackendURL != "" {
		if uc.isPotentialMismatch(cfg.AccessToken, cfg.BackendURL) {
			issues = append(issues, ValidationIssue{
				Key:        "config",
				Message:    "Access token and backend URL may not match",
				Severity:   ValidationSeverityWarning,
				Suggestion: "Verify that your token is valid for the specified backend URL",
			})
		}
	}

	return issues
}

// Helper methods

func (uc *validateConfigUseCase) isValidTokenFormat(token string) bool {
	// Basic token format validation
	// This could be enhanced based on your actual token format
	if len(token) < 20 {
		return false
	}

	// Check for alphanumeric characters (basic check)
	for _, char := range token {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_' || char == '.') {
			return false
		}
	}

	return true
}

func (uc *validateConfigUseCase) isValidURL(url string) bool {
	url = strings.TrimSpace(url)
	if len(url) == 0 {
		return false
	}

	// Check for basic URL structure
	hasProtocol := strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
	if !hasProtocol {
		return false
	}

	// Check for domain part after protocol
	if strings.HasPrefix(url, "http://") {
		domain := url[7:] // Remove "http://"
		return len(domain) > 0 && strings.Contains(domain, ".")
	}

	if strings.HasPrefix(url, "https://") {
		domain := url[8:] // Remove "https://"
		return len(domain) > 0 && strings.Contains(domain, ".")
	}

	return false
}

func (uc *validateConfigUseCase) testConnectivity(ctx context.Context, url string) *ValidationIssue {
	// Create a context with timeout for connectivity test
	testCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	// Create request
	req, err := http.NewRequestWithContext(testCtx, "HEAD", url, nil)
	if err != nil {
		return &ValidationIssue{
			Key:        "backend_url",
			Message:    "Failed to create connectivity test request",
			Severity:   ValidationSeverityWarning,
			Suggestion: "Check if the backend URL is accessible",
		}
	}

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return &ValidationIssue{
			Key:        "backend_url",
			Message:    "Backend URL is not reachable: " + err.Error(),
			Severity:   ValidationSeverityWarning,
			Suggestion: "Verify your internet connection and that the backend service is running",
		}
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 500 {
		return &ValidationIssue{
			Key:        "backend_url",
			Message:    "Backend service appears to be experiencing issues",
			Severity:   ValidationSeverityWarning,
			Suggestion: "The backend service may be temporarily unavailable",
		}
	}

	return nil // No connectivity issues
}

func (uc *validateConfigUseCase) isPotentialMismatch(token, url string) bool {
	// This is a placeholder for more sophisticated logic
	// You could implement checks like:
	// - Development tokens should only be used with development URLs
	// - Production tokens should only be used with production URLs
	// - Token issuer should match the backend domain

	// For now, just return false (no mismatch detected)
	return false
}
