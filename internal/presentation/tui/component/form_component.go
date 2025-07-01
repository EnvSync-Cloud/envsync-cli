package component

import (
	"github.com/charmbracelet/huh"
)

// FieldDef defines a single form field for the generic form component.
type FieldDef[T any] struct {
	Title       string
	Description string
	Placeholder string
	Value       *string
	Lines       int // For multi-line input, 0 means single-line
	Validate    func(string) error
}

// FormComponent builds and runs a huh form for any struct T with string fields.
type FormComponent[T any] struct {
	Fields []FieldDef[T]
	Theme  *huh.Theme
}

// NewFormComponent creates a new form component for the given fields.
func NewFormComponent[T any](fields []FieldDef[T], theme *huh.Theme) *FormComponent[T] {
	return &FormComponent[T]{
		Fields: fields,
		Theme:  theme,
	}
}

// Run displays the form and fills the struct fields.
func (fc *FormComponent[T]) Run() error {
	var groups []huh.Field

	for _, f := range fc.Fields {
		var input huh.Field
		if f.Lines > 1 {
			input = huh.NewText().
				Title(f.Title).
				Description(f.Description).
				Placeholder(f.Placeholder).
				Value(f.Value).
				Lines(f.Lines).
				Validate(f.Validate)
		} else {
			input = huh.NewInput().
				Title(f.Title).
				Description(f.Description).
				Placeholder(f.Placeholder).
				Value(f.Value).
				Validate(f.Validate)
		}
		groups = append(groups, input)
	}

	form := huh.NewForm(huh.NewGroup(groups...))
	if fc.Theme != nil {
		form = form.WithTheme(fc.Theme)
	}
	return form.Run()
}
