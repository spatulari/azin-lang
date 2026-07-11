package diagnostics

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/token"
)

type DiagnosticKind uint8

const (
	Note = iota
	Warning
	Error
)

type Diagnostic struct {
	Kind     DiagnosticKind
	Message  string
	Position token.Position
	Length   uint32
}

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

func (d Diagnostic) Error() string {
	return fmt.Sprintf("%s: %s", d.Kind, d.Message)
}
