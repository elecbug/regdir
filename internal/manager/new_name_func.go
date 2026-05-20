package manager

// RenameFunc defines a function type that takes FileWithDir and returns a string representing the new name for the file.
type RenameFunc func(FileWithDir) string
