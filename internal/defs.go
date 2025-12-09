// Package internal provides various internal functions and constants.
package internal

import "path/filepath"

const (
	// FILE_CONFIG is the configuration file name.
	FILE_CONFIG string = "config.json"
	// FILE_README is the README file name.
	FILE_README string = "HZMIND.md"
	// FILE_IGNORE is the ignore file name.
	FILE_IGNORE string = ".hzmignore"
)

const (
	// DIR_MAIN is the main directory name.
	DIR_MAIN string = "hzmind"
	// DIR_OUT is the output directory name.
	DIR_OUT string = "out"
)

var (
	// PATH_FILE_CONFIG is the full path to the configuration file.
	PATH_FILE_CONFIG string = FILE_CONFIG
	// PATH_FILE_README is the full path to the README file.
	PATH_FILE_README string = filepath.Join(DIR_MAIN, FILE_README)
	// PATH_FILE_IGNORE is the full path to the ignore file.
	PATH_FILE_IGNORE string = filepath.Join(DIR_MAIN, FILE_IGNORE)
	// PATH_DIR_OUT is the full path to the output directory.
	PATH_DIR_OUT string = filepath.Join(DIR_MAIN, DIR_OUT)
)
