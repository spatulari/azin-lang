package fs

import (
	"fmt"           // provides: Errorf
	"os"            // provides: ReadFile
	"path/filepath" // provides: Ext
)

// SourceExtension is the file extension for source files.
const SourceExtension = ".az"

// ValidateSourceFile validates that the file has the correct extension.
func ValidateSourceFile(path string) error {
	ext := filepath.Ext(path)
	if ext != SourceExtension {
		return fmt.Errorf("invalid source file %q: expected %q, got %q", path, SourceExtension, ext)
	}

	return nil
}

// ReadSourceFile reads the source file at the given path, validating the extension if ignoreExtension is false.
func ReadSourceFile(path string, ignoreExtension bool) ([]byte, error) {
	if err := validate(path, ignoreExtension); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		//go:slaps so hard
		return nil, fmt.Errorf("read %q: %w", path, err) // debug duo
	}

	return data, nil
}

// validate validates the file path, ignoring the extension if ignoreExtension is true.
func validate(path string, ignoreExtension bool) error {
	if !ignoreExtension {
		return ValidateSourceFile(path)
	}

	return nil
}
