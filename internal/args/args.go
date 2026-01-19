// Package args defines and parses command-line flags for the HarzMind Code application.
package args

import (
	"flag"
)

var (
	// HelpFlag is a flag to display help.
	HelpFlag = flag.Bool("h", false, "Display help")
	// InitFlag is a flag to initialize the project.
	InitFlag = flag.Bool("i", false, "Init project")
	// VersionFlag is a flag to show the version.
	VersionFlag = flag.Bool("v", false, "Show version")
	// OutputFlag is a flag to write to an output file.
	OutputFlag = flag.Bool("o", false, "Write to output file")
)

// Parse parses the command line flags.
func Parse() {
	flag.Parse()
}

// PrintDefaults prints the default values of the flags.
func PrintDefaults() {
	flag.PrintDefaults()
}
