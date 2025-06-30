package formatters

import (
	"fmt"
	"io"
	"strings"
)

type ConfigFormatter struct {
	*BaseFormatter
}

// NewConfigFormatter creates a new ConfigFormatter instance
func NewConfigFormatter() *ConfigFormatter {
	base := NewBaseFormatter()
	return &ConfigFormatter{
		BaseFormatter: base,
	}
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
