package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

const SourceExtension = ".az"

func ValidateSourceFile(path string) error {
	if ext := filepath.Ext(path); ext != SourceExtension {
		return fmt.Errorf(
			"invalid source file %q: expected %q extension, got %q",
			path,
			SourceExtension,
			ext,
		)
	}

	return nil
}

func ReadSourceFile(path string, ignoreExtension bool) ([]byte, error) {
	if !ignoreExtension {
		if err := ValidateSourceFile(path); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}

	return data, nil
}
