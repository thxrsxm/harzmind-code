package codebase

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Path    string `json:"path"`
}

func (f *File) String() string {
	return fmt.Sprintf("Name: %s, Content: %s, Path: %s", f.Name, f.Content, f.Path)
}

func genIgnoreFilter(files []string) string {
	var sb strings.Builder
	for i, v := range files {
		sb.WriteByte('(')
		sb.WriteString(v)
		sb.WriteByte(')')
		if i < len(files)-1 {
			sb.WriteByte('|')
		}
	}
	return sb.String()
}

func GetCodeBase(root string, ignoreFiles []string) []File {
	files := []File{}
	re := regexp.MustCompile(genIgnoreFilter(ignoreFiles))
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Check if the current entry matches the ignore filter
		if len(ignoreFiles) > 0 && re.MatchString(d.Name()) {
			if d.IsDir() {
				// Skip this directory and its contents
				return filepath.SkipDir
			} else {
				// Create and append File to slice without content
				files = append(files, File{
					Name:    d.Name(),
					Content: "",
					Path:    path,
				})
				return nil
			}
		}
		// If it's not a directory, process it as a file
		if !d.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("error opening file %s: %v", path, err)
			}
			defer file.Close()
			// Read file content
			content, err := io.ReadAll(file)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", path, err)
			}
			// Create and append File to slice
			files = append(files, File{
				Name:    d.Name(),
				Content: string(content),
				Path:    path,
			})
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err)
	}
	return files
}
