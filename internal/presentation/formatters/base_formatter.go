package formatters

import (
	"encoding/json"
	"fmt"
	"io"
)

type BaseFormatter struct{}

func NewBaseFormatter() *BaseFormatter {
	return &BaseFormatter{}
}

func (f *BaseFormatter) FormatJSON(writer io.Writer, data any) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal applications to JSON: %w", err)
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON output: %w", err)
	}

	// Add newline for better formatting
	_, err = writer.Write([]byte("\n"))
	return err
}
