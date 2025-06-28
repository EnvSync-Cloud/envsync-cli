package utils

import (
	"os"
)

// IsTerminal checks if the output is being written to a terminal
func IsTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// IsInteractiveMode determines if interactive mode should be used
func IsInteractiveMode(jsonFlag bool) bool {
	// Don't use interactive mode if:
	// 1. JSON output is requested
	// 2. Explicitly disabled with --no-interactive flag
	// 3. Not running in a terminal
	if jsonFlag || !IsTerminal() {
		return false
	}

	return true
}

// IsCI checks if running in a CI environment
func IsCI() bool {
	ci := os.Getenv("CI")
	return ci == "true" || ci == "1"
}

// GetTerminalSize returns the width and height of the terminal
func GetTerminalSize() (width, height int) {
	// Default fallback values
	width, height = 80, 24

	// You can add more sophisticated terminal size detection here
	// For now, returning defaults
	return width, height
}
