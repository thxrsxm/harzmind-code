package codebase

import (
	"sort"
	"strings"
)

// treeNode represents a node (directory or file) in a file tree structure.
type treeNode struct {
	name     string
	children map[string]*treeNode
}

// newTreeNode creates a new tree node.
func newTreeNode(name string) *treeNode {
	return &treeNode{
		name:     name,
		children: make(map[string]*treeNode),
	}
}

// Tree creates a tree structure of files and directories in the codebase
func Tree(files []File) string {
	// Create a map for the tree structure
	tree := make(map[string]*treeNode)
	// Add all files and directories to the tree structure
	for _, file := range files {
		// Normalize path separators to '/'
		normalizedPath := strings.ReplaceAll(file.Path, "\\", "/")
		// Extract the path and split it into components
		parts := strings.Split(normalizedPath, "/")
		current := tree
		// Iterate through all path components except the filename
		for _, part := range parts[:len(parts)-1] {
			// Skip empty parts (e.g., leading '/')
			if part == "" {
				continue
			}
			// Create a new node if it doesn't exist
			if _, exists := current[part]; !exists {
				current[part] = newTreeNode(part)
			}
			// Move to the next node
			current = current[part].children
		}
		// Add the file as a leaf node
		fileName := parts[len(parts)-1]
		if fileName != "" {
			current[fileName] = newTreeNode(fileName)
		}
	}
	// Build the tree structure as a string
	return buildTreeString(tree, "")
}

// buildTreeString recursively builds the tree structure as a string
func buildTreeString(nodes map[string]*treeNode, prefix string) string {
	var result strings.Builder
	// Sort the nodes alphabetically
	keys := make([]string, 0, len(nodes))
	for k := range nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// Iterate through all nodes
	for i, key := range keys {
		node := nodes[key]
		isLast := i == len(keys)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		// Add the current node
		result.WriteString(prefix + connector + key + "\n")
		// If the node has children, add them recursively
		if len(node.children) > 0 {
			newPrefix := prefix
			if !isLast {
				newPrefix += "│   "
			} else {
				newPrefix += "    "
			}
			result.WriteString(buildTreeString(node.children, newPrefix))
		}
	}
	return result.String()
}
