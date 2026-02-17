// Package args defines and parses command-line flags for the HarzMind Code application.
package args

import (
	"flag"
	"fmt"
	"os"
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
	// LogFlag is a flag to enable logging.
	LogFlag = flag.Bool("l", false, "Enable logging")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIMPORTANT:\n")
		fmt.Fprintf(os.Stderr, "   You must first initialize the project using the -i flag\n")
		fmt.Fprintf(os.Stderr, "   before using other features. Run '%s -i' to get started.\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -i           Initialize project\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -v           Show version\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -i -o -l     Init with output and logging\n", os.Args[0])
	}
}

// Parse parses the command line flags.
func Parse() {
	flag.Parse()
}

// PrintDefaults prints the default values of the flags.
func PrintDefaults() {
	flag.PrintDefaults()
}

func PrintUsage() {
	flag.Usage()
}
