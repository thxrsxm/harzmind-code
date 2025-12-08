package internal

import "path/filepath"

const (
	FILE_CONFIG string = "config.json"
	FILE_README string = "HZMIND.md"
	FILE_IGNORE string = ".hzmignore"
)

const (
	DIR_MAIN string = "hzmind"
	DIR_OUT  string = "out"
)

var (
	PATH_FILE_CONFIG string = FILE_CONFIG
	PATH_FILE_README string = filepath.Join(DIR_MAIN, FILE_README)
	PATH_FILE_IGNORE string = filepath.Join(DIR_MAIN, FILE_IGNORE)
	PATH_DIR_OUT     string = filepath.Join(DIR_MAIN, DIR_OUT)
)
