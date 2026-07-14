package codegen

import (
	"bytes"
	"path/filepath"

	"github.com/azin-lang/Azin/internal/ast"
)

// Transpiler the Transpiler struct is responsible for transpiling the AST to C code.
type Transpiler struct {
	buf    bytes.Buffer
	indent int
	enums  map[string]bool
}

// New create a new Transpiler.
func New() *Transpiler {
	return &Transpiler{
		enums: map[string]bool{},
	}
}

// handleImports handles the import statements in the AST and writes the corresponding C include statements to the buffer.
func (t *Transpiler) handleImports(program *ast.Program) string {
	for _, stmt := range program.Statements {
		if imp, ok := stmt.(*ast.ImportCStmt); ok {
			header := imp.Path.Value
			if filepath.Ext(header) == "" {
				header += ".h"
			}
			t.write("#include <" + header + ">\n")
		}
	}

	t.write("#include <stdbool.h>") // for bool, true, false because *great* C doesn't have those built-in

	t.newline()
	t.newline()

	return t.buf.String()
}

// Transpile takes an AST and transpiles it to C code.
func (t *Transpiler) Transpile(program *ast.Program) string {
	t.buf.Reset()

	t.handleImports(program)

	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.EnumStmt:
			t.compileEnum(s)
			t.newline()
		case *ast.StructStmt:
			t.compileStruct(s)
			t.newline()
		}
	}

	for _, stmt := range program.Statements {
		switch stmt.(type) {
		case *ast.ImportCStmt, *ast.StructStmt, *ast.EnumStmt:
			continue
		default:
			t.compileStatement(stmt)
			t.newline()
		}
	}

	return t.buf.String()
}
