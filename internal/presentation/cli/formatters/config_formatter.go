package formatters

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/EnvSync-Cloud/envsync-cli/internal/config"
)

type ConfigFormatter struct{}

func NewConfigFormatter() *ConfigFormatter {
	return &ConfigFormatter{}
}

// FormatJSON formats config as JSON
func (f *ConfigFormatter) FormatJSON(writer io.Writer, cfg config.AppConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write JSON output: %w", err)
	}

	// Add newline for better formatting
	_, err = writer.Write([]byte("\n"))
	return err
}

// FormatTable formats config as a readable table
func (f *ConfigFormatter) FormatTable(writer io.Writer, cfg config.AppConfig) error {
	header := "📋 Current Configuration:\n"
	if _, err := writer.Write([]byte(header)); err != nil {
		return err
	}

	// Access Token
	accessToken := cfg.AccessToken
	if accessToken == "" {
		accessToken = "<not set>"
	} else {
		// Mask the token for security
		accessToken = f.maskToken(accessToken)
	}
	line := fmt.Sprintf("🔑 access_token: %s\n", accessToken)
	if _, err := writer.Write([]byte(line)); err != nil {
		return err
	}

	// Backend URL
	backendURL := cfg.BackendURL
	if backendURL == "" {
		backendURL = "<not set>"
	}
	line = fmt.Sprintf("🌐 backend_url: %s\n", backendURL)
	if _, err := writer.Write([]byte(line)); err != nil {
		return err
	}

	return nil
}

// FormatSingleValue formats a single config value
func (f *ConfigFormatter) FormatSingleValue(writer io.Writer, key, value string) error {
	var output string

	switch strings.ToLower(key) {
	case "access_token", "accesstoken":
		maskedValue := value
		if value != "" {
			maskedValue = f.maskToken(value)
		} else {
			maskedValue = "<not set>"
		}
		output = fmt.Sprintf("🔑 access_token: %s\n", maskedValue)
	case "backend_url", "backendurl":
		if value == "" {
			value = "<not set>"
		}
		output = fmt.Sprintf("🌐 backend_url: %s\n", value)
	default:
		output = fmt.Sprintf("❓ %s: %s\n", key, value)
	}

	_, err := writer.Write([]byte(output))
	return err
}

// FormatCompact formats config in compact format
func (f *ConfigFormatter) FormatCompact(writer io.Writer, cfg config.AppConfig) error {
	var parts []string

	// Access token status
	if cfg.AccessToken != "" {
		parts = append(parts, "🔑 token: set")
	} else {
		parts = append(parts, "🔑 token: not set")
	}

	// Backend URL status
	if cfg.BackendURL != "" {
		parts = append(parts, fmt.Sprintf("🌐 url: %s", cfg.BackendURL))
	} else {
		parts = append(parts, "🌐 url: not set")
	}

	output := strings.Join(parts, " | ") + "\n"
	_, err := writer.Write([]byte(output))
	return err
}

// FormatValidationResult formats config validation results
func (f *ConfigFormatter) FormatValidationResult(writer io.Writer, isValid bool, issues []string) error {
	if isValid {
		output := "✅ Configuration is valid\n"
		_, err := writer.Write([]byte(output))
		return err
	}

	// Format validation issues
	header := "❌ Configuration validation failed:\n"
	if _, err := writer.Write([]byte(header)); err != nil {
		return err
	}

	for _, issue := range issues {
		line := fmt.Sprintf("   • %s\n", issue)
		if _, err := writer.Write([]byte(line)); err != nil {
			return err
		}
	}

	return nil
}

// Helper methods

func (f *ConfigFormatter) maskToken(token string) string {
	if len(token) <= 8 {
		return strings.Repeat("*", len(token))
	}

	// Show first 4 and last 4 characters
	prefix := token[:4]
	suffix := token[len(token)-4:]
	middle := strings.Repeat("*", len(token)-8)

	return prefix + middle + suffix
}

// FormatSuccess formats success messages
func (f *ConfigFormatter) FormatSuccess(writer io.Writer, message string) error {
	output := fmt.Sprintf("✅ %s\n", message)
	_, err := writer.Write([]byte(output))
	return err
}

// FormatError formats error messages
func (f *ConfigFormatter) FormatError(writer io.Writer, message string) error {
	output := fmt.Sprintf("❌ %s\n", message)
	_, err := writer.Write([]byte(output))
	return err
}

// FormatWarning formats warning messages
func (f *ConfigFormatter) FormatWarning(writer io.Writer, message string) error {
	output := fmt.Sprintf("⚠️  %s\n", message)
	_, err := writer.Write([]byte(output))
	return err
}

// FormatInfo formats info messages
func (f *ConfigFormatter) FormatInfo(writer io.Writer, message string) error {
	output := fmt.Sprintf("ℹ️  %s\n", message)
	_, err := writer.Write([]byte(output))
	return err
}

// FormatKeyValueList formats a list of key-value pairs
func (f *ConfigFormatter) FormatKeyValueList(writer io.Writer, title string, items map[string]string) error {
	if len(items) == 0 {
		return nil
	}

	header := fmt.Sprintf("📋 %s:\n", title)
	if _, err := writer.Write([]byte(header)); err != nil {
		return err
	}

	for key, value := range items {
		line := fmt.Sprintf("   • %s: %s\n", key, value)
		if _, err := writer.Write([]byte(line)); err != nil {
			return err
		}
	}

	return nil
}
