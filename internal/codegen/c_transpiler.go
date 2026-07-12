package codegen

import (
	"bytes"
	"path/filepath"

	"github.com/azin-lang/Azin/internal/ast"
)

// the Transpiler struct is responsible for transpiling the AST to C code.
type Transpiler struct {
	buf    bytes.Buffer
	indent int
}

// create a new Transpiler.
func New() *Transpiler {
	return &Transpiler{}
}

// Transpile transpiles the AST to C code.
func (t *Transpiler) Transpile(program *ast.Program) string {
	hasImports := false

	for _, stmt := range program.Statements {
		if imp, ok := stmt.(*ast.ImportCStmt); ok {
			header := imp.Path.Value

			if filepath.Ext(header) == "" {
				header += ".h"
			}

			t.printf("#include <%s>\n", header)
			hasImports = true
		}
	}

	if hasImports {
		t.newline()
	}

	for _, stmt := range program.Statements {
		if s, ok := stmt.(*ast.StructStmt); ok {
			t.compileStruct(s)
		}
	}

	for _, stmt := range program.Statements {
		if _, ok := stmt.(*ast.StructStmt); ok {
			continue
		}

		t.compileStatement(stmt)
		t.newline()
	}

	return t.buf.String()
}
