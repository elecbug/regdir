package manager

import (
	"fmt"
	"os"
	"path/filepath"
)

// ColletAllFiles collects all files in the specified directory that satisfy the given condition.
func ColletAllFiles(root string, condition ConditionFunc) (Files, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("provided path is not a directory")
	}

	files := make(Files, 0)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}

		if condition(info) {
			files = append(files, FileWithDir{FileInfo: info, ParentDir: filepath.Dir(path)})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to collect files: %w", err)
	}

	return files, nil
}
