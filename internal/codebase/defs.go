package codebase

import (
	"github.com/thxrsxm/harzmind-code/internal/common"
)

// ignorePatterns is a list of default patterns to ignore when retrieving the codebase.
var ignorePatterns []string = []string{
	".git",
	".idea",
	".vscode",
	"node_modules",
	"vendor",
	"*.exe",
	"config.xml",
	common.FILE_IGNORE,
	common.FILE_README,
	common.DIR_MAIN + "/",
}
