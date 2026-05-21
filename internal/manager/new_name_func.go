package manager

import (
	"strings"
)

// RenameFunc defines a function type that takes a FileWithDir and returns a new name for the file.
type RenameFunc func(FileWithDir) string

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
