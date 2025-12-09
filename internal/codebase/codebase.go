// Package codebase provides functionality for handling a codebase.
// It allows for retrieving a list of files within a given directory,
// excluding certain files and directories based on ignore patterns.
package codebase

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/thxrsxm/harzmind-code/internal"
)

// File represents a file within the codebase.
// It contains the file's name, content, and path.
type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Path    string `json:"path"`
}

// String returns a string representation of the File struct.
func (f File) String() string {
	return fmt.Sprintf("Name: %s, Content: %s, Path: %s", f.Name, f.Content, f.Path)
}

// createIgnorer creates a new ignore compiler based on the ignore patterns.
// It also reads additional ignore patterns from a .hzmignore file if it exists.
func createIgnorer() *ignore.GitIgnore {
	// Start with predefined ignore patterns
	patterns := make([]string, len(ignorePatterns))
	copy(patterns, ignorePatterns)
	// Check if .hzmignore file exists and add its patterns
	if file, err := os.Open(internal.PATH_FILE_IGNORE); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// Ignore empty lines and comments
			if line != "" && line[0] != '#' {
				patterns = append(patterns, line)
			}
		}
	}
	return ignore.CompileIgnoreLines(patterns...)
}

// IgnoreFileExists checks if the .hzmignore file exists.
func IgnoreFileExists() bool {
	return internal.FileExists(internal.PATH_FILE_IGNORE)
}

// GetCodeBase retrieves a list of files within the given root directory,
// excluding files and directories based on ignore patterns.
func GetCodeBase(root string) ([]File, error) {
	files := []File{}
	ignorer := createIgnorer()
	// Walk through the directory tree
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Error accessing file/directory
			return err
		}
		// Check if the current path should be ignored
		if ignorer.MatchesPath(path) {
			if d.IsDir() {
				// Skip entire directory and subdirectories
				return filepath.SkipDir
			}
			// Skip individual file
			return nil
		}
		// If we reach this point, the file/directory should not be ignored
		if !d.IsDir() {
			// Open the file
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("error opening file %s: %v", path, err)
			}
			defer file.Close()
			// Read the file content
			content, err := io.ReadAll(file)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", path, err)
			}
			// Add the file to the list
			files = append(files, File{
				Name:    d.Name(),
				Content: string(content),
				Path:    path,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}
	return files, nil
}
