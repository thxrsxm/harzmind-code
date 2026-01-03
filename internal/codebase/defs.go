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
