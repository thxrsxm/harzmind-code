// Package codebase provides functionality for handling a codebase.
// It allows for retrieving a list of files within a given directory,
// excluding certain files and directories based on ignore patterns.
package codebase

import "github.com/thxrsxm/harzmind-code/internal"

// ignorePatterns is a list of default patterns to ignore when retrieving the codebase.
var ignorePatterns []string = []string{
	".git",
	".idea",
	".vscode",
	"node_modules",
	"vendor",
	"*.exe",
	"config.xml",
	internal.FILE_IGNORE,
	internal.FILE_README,
	internal.DIR_MAIN + "/",
}
