package diagnostics

import (
	"fmt"     // provides: Println
	"strings" // provides: Builder
)

// Engine is responsible for managing diagnostics.
type Engine struct {
	diagnostics []Diagnostic
}

// New creates a new Engine.
func New() *Engine {
	return &Engine{}
}

// Report adds a diagnostic to the engine.
func (e *Engine) Report(kind DiagnosticKind, format string, args ...any) {
	e.diagnostics = append(e.diagnostics, Diagnostic{
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
	})
}

// ReportError adds an error diagnostic to the engine.
func (e *Engine) ReportError(format string, args ...any) {
	e.Report(Error, format, args...)
}

// ReportWarning adds a warning diagnostic to the engine.
func (e *Engine) ReportWarning(format string, args ...any) {
	e.Report(Warning, format, args...)
}

// ReportNote adds a note diagnostic to the engine.
func (e *Engine) ReportNote(format string, args ...any) {
	e.Report(Note, format, args...)
}

// Diagnostics returns all diagnostics.
func (e *Engine) Diagnostics() []Diagnostic {
	return e.diagnostics
}

// HasErrors returns true if the engine has errors.
func (e *Engine) HasErrors() bool {
	for _, d := range e.diagnostics {
		if d.Kind == Error {
			return true
		}
	}

	return false
}

// Err returns an error if the engine has errors.
func (e *Engine) Err() error {
	if !e.HasErrors() {
		return nil
	}
	return e
}

// Error returns the string representation of the engine's diagnostics.
func (e *Engine) Error() string {
	if len(e.diagnostics) == 0 {
		return ""
	}

	var b strings.Builder
	for i, d := range e.diagnostics {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(d.Error())
	}

	return b.String()
}
