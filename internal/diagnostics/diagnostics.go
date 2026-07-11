package diagnostics

import "fmt" // provides: Println

// DiagnosticKind represents the kind of diagnostic message.
type DiagnosticKind uint8

const (
	Note = iota
	Warning
	Error
)

// String returns the string representation of the diagnostic kind.
func (k DiagnosticKind) String() string {
	switch k {
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Note:
		return "note"
	default:
		return "unknown"
	}
}

// Diagnostic represents a diagnostic message.
type Diagnostic struct {
	Kind    DiagnosticKind
	Message string
	Line    int
	Column  int
	Offset  int
	Length  int
}

// Error returns the string representation of the diagnostic message.
func (d Diagnostic) Error() string {
	return fmt.Sprintf("%s: %s", d.Kind, d.Message)
}
