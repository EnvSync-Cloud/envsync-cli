package formatters

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/EnvSync-Cloud/envsync-cli/internal/domain"
)

type AppFormatter struct{}

func NewAppFormatter() *AppFormatter {
	return &AppFormatter{}
}

// FormatJSON formats applications as JSON
func (f *AppFormatter) FormatJSON(writer io.Writer, apps []domain.Application) error {
	data, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal applications to JSON: %w", err)
	}

	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write JSON output: %w", err)
	}

	// Add newline for better formatting
	_, err = writer.Write([]byte("\n"))
	return err
}

// FormatTable formats applications as a readable table
func (f *AppFormatter) FormatTable(writer io.Writer, apps []domain.Application) error {
	if len(apps) == 0 {
		_, err := writer.Write([]byte("📭 No applications found.\n"))
		return err
	}

	// Header
	header := "🚀 Available Applications:\n"
	if _, err := writer.Write([]byte(header)); err != nil {
		return err
	}

	// Applications
	for i, app := range apps {
		if i > 0 {
			// Add separator between apps
			separator := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
			if _, err := writer.Write([]byte(separator)); err != nil {
				return err
			}
		}

		appOutput := f.formatSingleApp(app)
		if _, err := writer.Write([]byte(appOutput)); err != nil {
			return err
		}
	}

	return nil
}

// FormatSingle formats a single application
func (f *AppFormatter) FormatSingle(writer io.Writer, app domain.Application) error {
	output := f.formatSingleApp(app)
	_, err := writer.Write([]byte(output))
	return err
}

// FormatList formats applications as a simple list
func (f *AppFormatter) FormatList(writer io.Writer, apps []domain.Application) error {
	if len(apps) == 0 {
		_, err := writer.Write([]byte("📭 No applications found.\n"))
		return err
	}

	for i, app := range apps {
		line := fmt.Sprintf("%d. 📛 %s (🆔 %s)\n", i+1, app.Name, app.ID)
		if _, err := writer.Write([]byte(line)); err != nil {
			return err
		}
	}

	return nil
}

// FormatCompact formats applications in compact format
func (f *AppFormatter) FormatCompact(writer io.Writer, apps []domain.Application) error {
	if len(apps) == 0 {
		_, err := writer.Write([]byte("📭 No applications found.\n"))
		return err
	}

	for _, app := range apps {
		line := fmt.Sprintf("📛 %s | 🆔 %s | 🌍 %s envs\n",
			app.Name,
			app.ID,
			f.getEnvCountDisplay(app.EnvCount))
		if _, err := writer.Write([]byte(line)); err != nil {
			return err
		}
	}

	return nil
}

// Helper methods

func (f *AppFormatter) formatSingleApp(app domain.Application) string {
	var output strings.Builder

	// Name
	output.WriteString(fmt.Sprintf("📛 Name: %s\n", app.Name))

	// ID
	output.WriteString(fmt.Sprintf("🆔 ID: %s\n", app.ID))

	// Description
	if app.Description != "" {
		output.WriteString(fmt.Sprintf("📝 Description: %s\n", app.Description))
	}

	// Organization ID
	if app.OrgID != "" {
		output.WriteString(fmt.Sprintf("🏢 Organization ID: %s\n", app.OrgID))
	}

	// Environment count
	if app.EnvCount != "" {
		envDisplay := f.getEnvCountDisplay(app.EnvCount)
		output.WriteString(fmt.Sprintf("🌍 Environments: %s\n", envDisplay))
	}

	// Environment types
	if len(app.EnvTypes) > 0 {
		output.WriteString("🏷️  Environment Types:\n")
		for _, envType := range app.EnvTypes {
			output.WriteString(fmt.Sprintf("   • %s (%s)\n", envType.Name, envType.ID))
		}
	}

	// Metadata
	if len(app.Metadata) > 0 {
		output.WriteString("🏷️  Metadata:\n")
		for key, value := range app.Metadata {
			output.WriteString(fmt.Sprintf("   • %s: %v\n", key, value))
		}
	}

	// Timestamps
	if !app.CreatedAt.IsZero() {
		output.WriteString(fmt.Sprintf("⏰ Created: %s\n", app.CreatedAt.Format("2006-01-02 15:04:05")))
	}

	if !app.UpdatedAt.IsZero() {
		output.WriteString(fmt.Sprintf("⏰ Updated: %s\n", app.UpdatedAt.Format("2006-01-02 15:04:05")))
	}

	return output.String()
}

func (f *AppFormatter) getEnvCountDisplay(envCount string) string {
	if envCount == "" {
		return "0"
	}
	return envCount
}

// FormatSuccess formats success messages
func (f *AppFormatter) FormatSuccess(writer io.Writer, message string) error {
	output := fmt.Sprintf("✅ %s\n", message)
	_, err := writer.Write([]byte(output))
	return err
}

// FormatError formats error messages
func (f *AppFormatter) FormatError(writer io.Writer, message string) error {
	output := fmt.Sprintf("❌ %s\n", message)
	_, err := writer.Write([]byte(output))
	return err
}

// FormatWarning formats warning messages
func (f *AppFormatter) FormatWarning(writer io.Writer, message string) error {
	output := fmt.Sprintf("⚠️  %s\n", message)
	_, err := writer.Write([]byte(output))
	return err
}

// FormatInfo formats info messages
func (f *AppFormatter) FormatInfo(writer io.Writer, message string) error {
	output := fmt.Sprintf("ℹ️  %s\n", message)
	_, err := writer.Write([]byte(output))
	return err
}
