package manager

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ConditionFunc defines a function type that takes os.FileInfo and returns a boolean.
type ConditionFunc func(os.FileInfo) bool

// WithExtension returns a ConditionFunc that checks if a file has the specified extension.
func WithExtension(ext string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && filepath.Ext(info.Name()) == ext
	}
}

// WithoutExtension returns a ConditionFunc that checks if a file does not have the specified extension.
func WithoutExtension(ext string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && filepath.Ext(info.Name()) != ext
	}
}

// WithPrefix returns a ConditionFunc that checks if a file name has the specified prefix.
func WithPrefix(prefix string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && strings.HasPrefix(info.Name(), prefix)
	}
}

// WithoutPrefix returns a ConditionFunc that checks if a file name does not have the specified prefix.
func WithoutPrefix(prefix string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && !strings.HasPrefix(info.Name(), prefix)
	}
}

// WithSuffix returns a ConditionFunc that checks if a file name has the specified suffix.
func WithSuffix(suffix string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && strings.HasSuffix(info.Name(), suffix)
	}
}

// WithoutSuffix returns a ConditionFunc that checks if a file name does not have the specified suffix.
func WithoutSuffix(suffix string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && !strings.HasSuffix(info.Name(), suffix)
	}
}

// WithSubstring returns a ConditionFunc that checks if a file name contains the specified substring.
func WithSubstring(substring string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && strings.Contains(info.Name(), substring)
	}
}

// WithoutSubstring returns a ConditionFunc that checks if a file name does not contain the specified substring.
func WithoutSubstring(substring string) ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && !strings.Contains(info.Name(), substring)
	}
}

// IsHidden returns a ConditionFunc that checks if a file is hidden (starts with a dot).
func IsHidden() ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && strings.HasPrefix(info.Name(), ".")
	}
}

// IsNotHidden returns a ConditionFunc that checks if a file is not hidden (does not start with a dot).
func IsNotHidden() ConditionFunc {
	return func(info os.FileInfo) bool {
		return !info.IsDir() && !strings.HasPrefix(info.Name(), ".")
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

// FindWithWildcard returns a ConditionFunc that checks if a file name matches the specified wildcard pattern.
func FindWithWildcard(pattern string, reverse bool) ConditionFunc {
	return func(info os.FileInfo) bool {
		matched, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return false
		}

		if reverse {
			return !matched
		} else {
			return matched
		}
	}
}

// FindWithRegex returns a ConditionFunc that checks if a file name matches the specified regular expression pattern.
func FindWithRegex(pattern string, reverse bool) ConditionFunc {
	return func(info os.FileInfo) bool {
		matched, err := regexp.MatchString(pattern, info.Name())
		if err != nil {
			return false
		}
		if reverse {
			return !matched
		} else {
			return matched
		}
	}
}
