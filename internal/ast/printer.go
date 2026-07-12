package ast

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PrintJSON serializes the AST node to JSON and prints it to standard output.
func PrintJSON(node Node) error {
	data, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal AST: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

// ExportJSON serializes the AST node to JSON and writes it to the specified destination path.
// The destPath should include the desired file name (e.g., "out/ast.json").
func ExportJSON(node Node, destPath string) error {
	data, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal AST: %w", err)
	}

	// Ensure the parent directories exist based on the provided file path
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write AST file: %w", err)
	}

	return nil
}
