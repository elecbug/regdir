package manager_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/elecbug/regdir/internal/manager"
)

// TestColletAllFiles tests the ColletAllFiles function with various conditions.
func TestColletAllFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// t.Logf("Temporary directory created: %s", tempDir)

	// Create some test files and directories
	testFiles := []string{"file1.txt", "file2.txt", "file3.log", "test1.txt", "subdir/file4.txt", "subdir/test2.log"}
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Define a condition function to filter .txt files
	conditions := []manager.ConditionFunc{
		func(info os.FileInfo) bool {
			return !info.IsDir() && filepath.Ext(info.Name()) == ".txt"
		},
		func(info os.FileInfo) bool {
			return !info.IsDir() && filepath.Ext(info.Name()) == ".log"
		},
		func(info os.FileInfo) bool {
			return !info.IsDir() && strings.HasPrefix(info.Name(), "file")
		},
		func(info os.FileInfo) bool {
			return info.IsDir()
		},
	}

	expectedCounts := []int{4, 2, 4, 2}
	expectedNames := []map[string]bool{
		{"file1.txt": true, "file2.txt": true, "test1.txt": true, "file4.txt": true},
		{"file3.log": true, "test2.log": true},
		{"file1.txt": true, "file2.txt": true, "file3.log": true, "file4.txt": true},
		{filepath.Base(tempDir): true, "subdir": true},
	}

	for i := range conditions {
		// Call the function under test
		files, err := manager.ColletAllFiles(tempDir, conditions[i])
		if err != nil {
			t.Fatalf("ColletAllFiles returned an error: %v", err)
		}

		// Verify the results
		if len(files) != expectedCounts[i] {
			t.Errorf("Expected %d files, got %d", expectedCounts[i], len(files))
		}

		for _, fileInfo := range files {
			if !expectedNames[i][fileInfo.Name()] {
				t.Errorf("Unexpected file found: %s", fileInfo.Name())
			}
		}
	}
}

// TestChangeFileNames tests the ChangeFileNames function to ensure it correctly renames files based on the provided newNameFunc.
func TestChangeFileNames(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// t.Logf("Temporary directory created: %s", tempDir)

	// Create some test files and directories
	testFiles := []string{"file1.txt", "file2.txt", "file3.log", "test1.txt", "subdir/file4.txt", "subdir/test2.log"}
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Define a condition function to filter .txt files
	conditions := []manager.ConditionFunc{
		func(info os.FileInfo) bool {
			return !info.IsDir() && strings.HasPrefix(info.Name(), "file")
		},
		func(info os.FileInfo) bool {
			return !info.IsDir() && filepath.Ext(info.Name()) == ".txt"
		},
		func(info os.FileInfo) bool {
			return !info.IsDir() && filepath.Ext(info.Name()) == ".log"
		},
		func(info os.FileInfo) bool {
			return info.IsDir() && info.Name() != filepath.Base(tempDir)
		},
	}

	resultsFiles := []map[string]bool{
		{"new_file1.txt": true, "new_file2.txt": true, "new_file3.log": true, "new_file4.txt": true},
		{"new_new_file1.txt": true, "new_new_file2.txt": true, "new_test1.txt": true, "new_new_file4.txt": true},
		{"new_new_file3.log": true, "new_test2.log": true},
		{"new_subdir": true},
	}

	for i := range conditions {
		// Call the function under test
		files, err := manager.ColletAllFiles(tempDir, conditions[i])
		if err != nil {
			t.Fatalf("ColletAllFiles returned an error: %v", err)
		}

		newNameFunc := func(file manager.FileWithDir) string {
			return "new_" + file.Name()
		}

		err = manager.ChangeFileNames(files, newNameFunc, false)
		if err != nil {
			t.Fatalf("ChangeFileNames returned an error: %v", err)
		}

		// Verify the results
		for _, fileInfo := range files {
			newName := "new_" + fileInfo.Name()
			if !resultsFiles[i][newName] {
				t.Errorf("Unexpected file name after change: %s", newName)
			}
		}
	}

	testFiles = []string{"file1.txt", "file2.txt"}
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Test overwrite scenario
	files, err := manager.ColletAllFiles(tempDir, func(info os.FileInfo) bool {
		return !info.IsDir() && info.Name() == "file1.txt"
	})
	if err != nil {
		t.Fatalf("ColletAllFiles returned an error: %v", err)
	}

	newNameFunc := func(file manager.FileWithDir) string {
		return "file2.txt" // This will cause a collision with an existing file
	}

	err = manager.ChangeFileNames(files, newNameFunc, false)
	if err == nil {
		t.Fatalf("Expected error due to name collision, but got nil")
	}

	if !strings.Contains(err.Error(), "target already exists") {
		t.Errorf("Expected error message to contain 'target already exists', got: %v", err)
	}

	err = manager.ChangeFileNames(files, newNameFunc, true)
	if err != nil {
		t.Fatalf("ChangeFileNames with overwrite returned an error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "file2.txt")); err != nil {
		t.Errorf("Expected file2.txt to exist after overwrite, but got error: %v", err)
	}
}

// TestChangeFileNamesWithSecondRules tests the ChangeFileNamesWithSecondRules function to ensure it correctly renames files based on the provided newNameFunc and secondRuleFuncs.
func TestChangeFileNamesWithSecondRules(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// t.Logf("Temporary directory created: %s", tempDir)

	filesToCreate := []string{
		"a.txt",
		"new_a.txt",
		"b.txt",
	}

	for _, name := range filesToCreate {
		path := filepath.Join(tempDir, name)
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	condition := func(info os.FileInfo) bool {
		return !info.IsDir() && info.Name() == "a.txt"
	}

	files, err := manager.ColletAllFiles(tempDir, condition)
	if err != nil {
		t.Fatalf("ColletAllFiles returned an error: %v", err)
	}

	firstRule := func(file manager.FileWithDir) string {
		return "new_" + file.Name()
	}

	secondRules := []manager.RenameFunc{
		func(file manager.FileWithDir) string {
			return "second_" + file.Name()
		},
	}

	err = manager.ChangeFileNamesWithSecondRules(files, firstRule, secondRules, false)
	if err != nil {
		t.Fatalf("ChangeFileNamesWithSecondRules returned an error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "second_a.txt")); err != nil {
		t.Errorf("expected second rule result file to exist: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "a.txt")); !os.IsNotExist(err) {
		t.Errorf("expected original file to be renamed")
	}

	if _, err := os.Stat(filepath.Join(tempDir, "new_a.txt")); err != nil {
		t.Errorf("expected pre-existing collision file to remain: %v", err)
	}
}

// TestChangeFileNamesWithPattern tests the ChangeFileNamesWithPattern function to ensure it correctly renames files based on the provided pattern and pattern type.
func TestChangeFileNamesWithPattern(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// t.Logf("Temporary directory created: %s", tempDir)

	// Create some test files and directories
	testFiles := []string{"file1.txt", "file2.txt", "file3.log", "test1.txt", "subdir/file4.txt", "subdir/test2.log"}
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	err = manager.ChangeFileNamesWithPattern(tempDir, "f*.txt", "new_f*.txt", manager.WILDCARD, false)
	if err != nil {
		t.Fatalf("ChangeFileNamesWithPattern returned an error: %v", err)
	}

	expectedFiles := []string{"new_file1.txt", "new_file2.txt", "file3.log", "test1.txt", "subdir/new_file4.txt", "subdir/test2.log"}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(filepath.Join(tempDir, expectedFile)); err != nil {
			t.Errorf("Expected file %s to exist, but got error: %v", expectedFile, err)
		}
	}

	err = manager.ChangeFileNamesWithPattern(tempDir, `new_file(\d+)\.txt`, `f$1.txt`, manager.REGEX, false)
	if err != nil {
		t.Fatalf("ChangeFileNamesWithPattern returned an error: %v", err)
	}

	expectedFiles = []string{"f1.txt", "f2.txt", "file3.log", "test1.txt", "subdir/f4.txt", "subdir/test2.log"}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(filepath.Join(tempDir, expectedFile)); err != nil {
			t.Errorf("Expected file %s to exist, but got error: %v", expectedFile, err)
		}
	}
}

// TestMoveFiles tests the MoveFiles function to ensure it correctly moves files to the target directory, handling overwrite scenarios as expected.
func TestMoveFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// t.Logf("Temporary directory created: %s", tempDir)

	// Create a target directory
	targetDir := filepath.Join(tempDir, "target")
	if err := os.Mkdir(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target directory: %v", err)
	}

	// Create some test files
	testFiles := []string{"file1.txt", "file2.txt"}
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	filesToMove, err := manager.ColletAllFiles(tempDir, func(info os.FileInfo) bool {
		return !info.IsDir() && (info.Name() == "file1.txt" || info.Name() == "file2.txt")
	})
	if err != nil {
		t.Fatalf("ColletAllFiles returned an error: %v", err)
	}

	err = manager.MoveFiles(filesToMove, targetDir, false)
	if err != nil {
		t.Fatalf("MoveFiles returned an error: %v", err)
	}

	for _, fileName := range testFiles {
		if _, err := os.Stat(filepath.Join(targetDir, fileName)); err != nil {
			t.Errorf("Expected file %s to be moved to target directory, but got error: %v", fileName, err)
		}
		if _, err := os.Stat(filepath.Join(tempDir, fileName)); !os.IsNotExist(err) {
			t.Errorf("Expected original file %s to be moved, but it still exists in the original location", fileName)
		}
	}

	// Test overwrite scenario
	if err := os.WriteFile(filepath.Join(targetDir, "file1.txt"), []byte("existing content"), 0644); err != nil {
		t.Fatalf("Failed to create existing file in target directory: %v", err)
	}

	filesToMove, err = manager.ColletAllFiles(tempDir, func(info os.FileInfo) bool {
		return !info.IsDir() && info.Name() == "file1.txt"
	})
	if err != nil {
		t.Fatalf("ColletAllFiles returned an error: %v", err)
	}

	err = manager.MoveFiles(filesToMove, targetDir, false)
	if err == nil {
		t.Fatalf("Expected error due to existing file in target directory, but got nil")
	}

	if !strings.Contains(err.Error(), "target already exists") {
		t.Errorf("Expected error message to contain 'target already exists', got: %v", err)
	}

	err = manager.MoveFiles(filesToMove, targetDir, true)
	if err != nil {
		t.Fatalf("MoveFiles with overwrite returned an error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "file1.txt")); err != nil {
		t.Errorf("Expected file1.txt to exist in target directory after overwrite, but got error: %v", err)
	}
}

// TestCopyFiles tests the CopyFiles function to ensure it correctly copies files to the target directory, handling overwrite scenarios as expected.
func TestCopyFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// t.Logf("Temporary directory created: %s", tempDir)

	// Create a target directory
	targetDir := filepath.Join(tempDir, "target")
	if err := os.Mkdir(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target directory: %v", err)
	}

	// Create some test files
	testFiles := []string{"file1.txt", "file2.txt"}
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	filesToCopy, err := manager.ColletAllFiles(tempDir, func(info os.FileInfo) bool {
		return !info.IsDir() && (info.Name() == "file1.txt" || info.Name() == "file2.txt")
	})
	if err != nil {
		t.Fatalf("ColletAllFiles returned an error: %v", err)
	}

	err = manager.CopyFiles(filesToCopy, targetDir, false)
	if err != nil {
		t.Fatalf("CopyFiles returned an error: %v", err)
	}

	for _, fileName := range testFiles {
		if _, err := os.Stat(filepath.Join(targetDir, fileName)); err != nil {
			t.Errorf("Expected file %s to be copied to target directory, but got error: %v", fileName, err)
		}
		if _, err := os.Stat(filepath.Join(tempDir, fileName)); err != nil {
			t.Errorf("Expected original file %s to still exist in original location, but got error: %v", fileName, err)
		}
	}

	// Test overwrite scenario
	if err := os.WriteFile(filepath.Join(targetDir, "file1.txt"), []byte("existing content"), 0644); err != nil {
		t.Fatalf("Failed to create existing file in target directory: %v", err)
	}

	filesToCopy, err = manager.ColletAllFiles(tempDir, func(info os.FileInfo) bool {
		return !info.IsDir() && info.Name() == "file1.txt"
	})
	if err != nil {
		t.Fatalf("ColletAllFiles returned an error: %v", err)
	}

	err = manager.CopyFiles(filesToCopy, targetDir, false)
	if err == nil {
		t.Fatalf("Expected error due to existing file in target directory, but got nil")
	}

	if !strings.Contains(err.Error(), "target already exists") {
		t.Errorf("Expected error message to contain 'target already exists', got: %v", err)
	}

	err = manager.CopyFiles(filesToCopy, targetDir, true)
	if err != nil {
		t.Fatalf("CopyFiles with overwrite returned an error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "file1.txt")); err != nil {
		t.Errorf("Expected file1.txt to exist in target directory after overwrite, but got error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tempDir, "file1.txt")); err != nil {
		t.Errorf("Expected original file1.txt to still exist in original location after copy, but got error: %v", err)
	}
}
