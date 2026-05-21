package manager

import (
	"regexp"
	"strings"
)

// RenameFunc defines a function type that takes a FileWithDir and returns a new name for the file.
type RenameFunc func(FileWithDir) string

// replaceWithWildcard replaces parts of the file name based on the oldPattern and newPattern, where '*' is used as a wildcard.
func replaceWithWildcard(file FileWithDir, oldPattern, newPattern string) string {
	name := file.Name()

	oldParts := strings.Split(oldPattern, "*")
	newParts := strings.Split(newPattern, "*")

	if len(oldParts) == 1 {
		if name == oldPattern {
			return newPattern
		}
		return name
	}

	if len(oldParts) != 2 || len(newParts) != 2 {
		return name
	}

	prefix := oldParts[0]
	suffix := oldParts[1]

	if !strings.HasPrefix(name, prefix) {
		return name
	}
	if !strings.HasSuffix(name, suffix) {
		return name
	}

	mid := name[len(prefix) : len(name)-len(suffix)]

	return newParts[0] + mid + newParts[1]
}

// replaceWithRegex is a placeholder function for future implementation of regex-based renaming.
// e.g. oldPattern: `file_(\d+)\.txt`, newPattern: `document_$1.txt` -> file_123.txt -> document_123.txt
func replaceWithRegex(file FileWithDir, oldPattern, newPattern string) string {
	name := file.Name()

	re, err := regexp.Compile(oldPattern)
	if err != nil {
		return name
	}

	return re.ReplaceAllString(name, newPattern)
}
