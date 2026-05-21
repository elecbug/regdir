package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ChangeFileNames changes the names of the files in the provided slice using the renameFunc to generate new names.
// If overwrite is false, it returns an error when the target path already exists.
func ChangeFileNames(files Files, renameFunc RenameFunc, overwrite bool) error {
	sort.Slice(files, func(i, j int) bool {
		return strings.Count(files[i].Path(), string(os.PathSeparator)) >
			strings.Count(files[j].Path(), string(os.PathSeparator))
	})

	var df = func(file FileWithDir) string {
		if file.IsDir() {
			return "directory"
		}
		return "file"
	}

	for _, file := range files {
		newName := renameFunc(file)
		if newName == "" {
			return fmt.Errorf("new name cannot be empty")
		}

		if newName != filepath.Base(newName) {
			return fmt.Errorf("new name must not contain path separators: %s", newName)
		}

		oldPath := filepath.Join(file.ParentDir, file.Name())
		newPath := filepath.Join(file.ParentDir, newName)

		if !overwrite {
			if _, err := os.Stat(newPath); err == nil {
				return fmt.Errorf("target already exists: %s", newPath)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("failed to check target path %s: %w", newPath, err)
			}
		}

		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("failed to rename %s: %w", df(file), err)
		}
	}

	return nil
}

// ChangeFileNamesWithSecondRules changes the names of the files in the provided slice using the renameFunc to generate new names.
func ChangeFileNamesWithSecondRules(files Files, renameFunc RenameFunc, secondRuleFuncs []RenameFunc, overwrite bool) error {
	sort.Slice(files, func(i, j int) bool {
		return strings.Count(files[i].Path(), string(os.PathSeparator)) >
			strings.Count(files[j].Path(), string(os.PathSeparator))
	})

	var df = func(file FileWithDir) string {
		if file.IsDir() {
			return "directory"
		}
		return "file"
	}

	var tryRename = func(file FileWithDir, oldPath string, newName string) error {
		if newName == "" {
			return fmt.Errorf("new name cannot be empty")
		}

		if newName != filepath.Base(newName) {
			return fmt.Errorf("new name must not contain path separators: %s", newName)
		}

		newPath := filepath.Join(file.ParentDir, newName)

		if !overwrite {
			if _, err := os.Stat(newPath); err == nil {
				return fmt.Errorf("target already exists: %s", newPath)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("failed to check target path %s: %w", newPath, err)
			}
		}

		return os.Rename(oldPath, newPath)
	}

	for _, file := range files {
		oldPath := filepath.Join(file.ParentDir, file.Name())

		err := tryRename(file, oldPath, renameFunc(file))
		if err == nil {
			continue
		}

		success := false
		lastErr := err

		for _, secondRuleFunc := range secondRuleFuncs {
			err = tryRename(file, oldPath, secondRuleFunc(file))
			if err == nil {
				success = true
				break
			}
			lastErr = err
		}

		if !success {
			return fmt.Errorf("failed to rename %s: %w", df(file), lastErr)
		}
	}

	return nil
}

// ChangeFileNamesWithPattern changes the names of files in the specified directory that match the oldPattern to newPattern.
// e.g. Wildcard, with oldPattern "file_*.txt" and newPattern "document_*.txt", "file_123.txt" will be renamed to "document_123.txt".
// e.g. Regex, with oldPattern `file_(\d+)\.txt` and newPattern `document_$1.txt`, "file_123.txt" will be renamed to "document_123.txt".
func ChangeFileNamesWithPattern(root string, oldPattern, newPattern string, patternType PatternType, overwrite bool) error {
	if oldPattern == newPattern {
		return fmt.Errorf("old pattern and new pattern are the same: %s", oldPattern)
	}

	var checkFunc ConditionFunc
	switch patternType {
	case REGEX:
		checkFunc = FindWithRegex(oldPattern, false)
	case WILDCARD:
		checkFunc = FindWithWildcard(oldPattern, false)
	default:
		return fmt.Errorf("unsupported pattern type: %v", patternType)
	}

	files, err := ColletAllFiles(root, checkFunc)
	if err != nil {
		return fmt.Errorf("failed to collect files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found matching the pattern: %s", oldPattern)
	}

	err = ChangeFileNames(files, func(file FileWithDir) string {
		switch patternType {
		case REGEX:
			return replaceWithRegex(file, oldPattern, newPattern)
		case WILDCARD:
			return replaceWithWildcard(file, oldPattern, newPattern)
		default:
			return file.Name()
		}
	}, overwrite)

	if err != nil {
		return fmt.Errorf("failed to change file names: %w", err)
	}

	return nil
}

// MoveFiles moves the specified files to the target directory. If overwrite is false,
// it returns an error when the target path already exists.
func MoveFiles(files Files, targetDir string, overwrite bool) error {
	for _, file := range files {
		oldPath := filepath.Join(file.ParentDir, file.Name())
		newPath := filepath.Join(targetDir, file.Name())

		if !overwrite {
			if _, err := os.Stat(newPath); err == nil {
				return fmt.Errorf("target already exists: %s", newPath)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("failed to check target path %s: %w", newPath, err)
			}
		}

		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("failed to move %s: %w", file.Name(), err)
		}
	}

	return nil
}

// CopyFiles copies the specified files to the target directory. If overwrite is false,
// it returns an error when the target path already exists.
func CopyFiles(files Files, targetDir string, overwrite bool) error {
	for _, file := range files {
		oldPath := filepath.Join(file.ParentDir, file.Name())
		newPath := filepath.Join(targetDir, file.Name())

		if !overwrite {
			if _, err := os.Stat(newPath); err == nil {
				return fmt.Errorf("target already exists: %s", newPath)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("failed to check target path %s: %w", newPath, err)
			}
		}

		input, err := os.ReadFile(oldPath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", oldPath, err)
		}

		err = os.WriteFile(newPath, input, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", newPath, err)
		}
	}

	return nil
}
