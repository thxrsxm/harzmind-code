package codebase

import "fmt"

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
