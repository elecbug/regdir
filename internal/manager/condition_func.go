package manager

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ConditionFunc defines a function type that takes os.FileInfo and returns a boolean.
type ConditionFunc func(os.FileInfo) bool

// HasExtension returns a ConditionFunc that checks if a file has the specified extension.
func HasExtension(ext string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && filepath.Ext(info.Name()) == ext
	}
}

// DoesNotHaveExtension returns a ConditionFunc that checks if a file does not have the specified extension.
func DoesNotHaveExtension(ext string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && filepath.Ext(info.Name()) != ext
	}
}

// HasPrefix returns a ConditionFunc that checks if a file name has the specified prefix.
func HasPrefix(prefix string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && strings.HasPrefix(info.Name(), prefix)
	}
}

// HasSuffix returns a ConditionFunc that checks if a file name has the specified suffix.
func HasSuffix(suffix string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && strings.HasSuffix(info.Name(), suffix)
	}
}

// HasSubstring returns a ConditionFunc that checks if a file name contains the specified substring.
func HasSubstring(substring string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && strings.Contains(info.Name(), substring)
	}
}

// IsDirectory returns a ConditionFunc that checks if the file is a directory.
func IsDirectory() ConditionFunc {
	return func(info os.FileInfo) bool {
		return info.IsDir()
	}
}

// IsFile returns a ConditionFunc that checks if the file is not a directory.
func IsFile() ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir()
	}
}

// ModifyTimeBefore returns a ConditionFunc that checks if a file's modification time is before the specified duration.
func ModifyTimeBefore(t time.Duration) ConditionFunc {
	return func(info os.FileInfo) bool {
		return time.Since(info.ModTime()) > t
	}
}

// ModifyTimeAfter returns a ConditionFunc that checks if a file's modification time is after the specified duration.
func ModifyTimeAfter(t time.Duration) ConditionFunc {
	return func(info os.FileInfo) bool {
		return time.Since(info.ModTime()) < t
	}
}
