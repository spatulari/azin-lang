package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/codegen"
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/parser"
	"github.com/azin-lang/Azin/internal/source"
)

func Compile(file *source.File, outputPath string) error {
	program, err := parseSource(file)
	if err != nil {
		return err
	}

	cCode := transpileToC(program)
	exeName := resolveExeName(outputPath)

	tmpPath, err := writeToTempFile(cCode)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	return runCompiler(tmpPath, exeName)
}

func parseSource(file *source.File) (*ast.Program, error) {
	diag := diagnostics.New(file)
	tokens := lexer.New(file, diag).Tokenize()
	if err := diag.Err(); err != nil {
		return nil, err
	}

	sourceString := string(file.Slice(0, file.Len()))
	p := parser.New(sourceString, tokens)
	program := p.ParseProgram()
	return program, diag.Err()
}

func transpileToC(program *ast.Program) string {
	tx := codegen.New()
	return tx.Transpile(program)
}

func resolveExeName(outputPath string) string {
	if outputPath == "" {
		return "output.exe"
	}
	if strings.HasSuffix(outputPath, ".c") {
		return strings.TrimSuffix(outputPath, ".c") + ".exe"
	}
	return outputPath
}

func writeToTempFile(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "azin_*.c")
	if err != nil {
		return "", fmt.Errorf("failed to create translation buffer: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to populate compile buffer: %w", err)
	}
	return tmpFile.Name(), nil
}

func runCompiler(sourcePath, exeName string) error {
	cmd := exec.Command("cl.exe", "/nologo", "/O2", "/Fe:"+exeName, sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("[MSVC] Compiling memory buffer via transient path %s...\n", filepath.Base(sourcePath))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cl.exe compilation failed: %w", err)
	}

	fmt.Printf("[Success] Executable generated directly to binary file: %s\n", exeName)
	return nil
}
