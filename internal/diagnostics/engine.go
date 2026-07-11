package diagnostics

import (
	"fmt"
	"strings"

	"github.com/azin-lang/Azin/internal/source"
	"github.com/azin-lang/Azin/internal/token"
)

type Engine struct {
	file        *source.File
	diagnostics []Diagnostic
	hasErrors   bool
}

func New(file *source.File) *Engine {
	return &Engine{
		file: file,
	}
}

func (e *Engine) Report(
	kind DiagnosticKind,
	pos token.Position,
	length uint32,
	format string,
	args ...any,
) {
	if kind == Error {
		e.hasErrors = true
	}

	e.diagnostics = append(e.diagnostics, Diagnostic{
		Kind:     kind,
		Message:  fmt.Sprintf(format, args...),
		Position: pos,
		Length:   length,
	})
}

func (e *Engine) ReportError(pos token.Position, length uint32, format string, args ...any) {
	e.Report(Error, pos, length, format, args...)
}

func (e *Engine) ReportWarning(pos token.Position, length uint32, format string, args ...any) {
	e.Report(Warning, pos, length, format, args...)
}

func (e *Engine) ReportNote(pos token.Position, length uint32, format string, args ...any) {
	e.Report(Note, pos, length, format, args...)
}

func (e *Engine) Diagnostics() []Diagnostic {
	return e.diagnostics
}

func (e *Engine) HasErrors() bool {
	return e.hasErrors
}

func (e *Engine) Err() error {
	if !e.HasErrors() {
		return nil
	}
	return e
}

func (e *Engine) LineColumn(pos token.Position) (line, column uint32) {
	return e.file.LineColumn(pos.Offset)
}

func (e *Engine) Text(tok token.Token) []byte {
	return e.file.Text(tok)
}

func (e *Engine) Line(line uint32) []byte {
	return e.file.Line(line)
}

func (e *Engine) Error() string {
	var b strings.Builder

	for i, d := range e.diagnostics {
		if i > 0 {
			b.WriteString("\n\n")
		}

		line, column := e.file.LineColumn(d.Position.Offset)

		fmt.Fprintf(
			&b,
			"%s:%d:%d: %s: %s\n",
			e.file.Name(),
			line,
			column,
			d.Kind,
			d.Message,
		)

		src := e.file.Line(line)
		b.Write(src)
		b.WriteByte('\n')

		prefix := src[:column-1]

		for _, ch := range prefix {
			if ch == '\t' {
				b.WriteByte('\t')
			} else {
				b.WriteByte(' ')
			}
		}

		b.WriteByte('^')
		for i := uint32(1); i < d.Length; i++ {
			b.WriteByte('~')
		}
	}

	return b.String()
}
