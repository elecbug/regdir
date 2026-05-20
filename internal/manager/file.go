package manager

import (
	"os"
	"path/filepath"
)

// FileWithDir is a struct that embeds os.FileInfo and includes an additional field for the parent directory.
type FileWithDir struct {
	os.FileInfo
	ParentDir string
}

// Path returns the full path of the file by joining the parent directory and the file name.
func (f FileWithDir) Path() string {
	return filepath.Join(f.ParentDir, f.Name())
}

// Files is a slice of FileWithDir, representing a collection of files.
type Files []FileWithDir
