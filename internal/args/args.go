package args

import (
	"flag"
)

var (
	HelpFlag    = flag.Bool("h", false, "Display help")
	InitFlag    = flag.Bool("i", false, "Init project")
	VersionFlag = flag.Bool("v", false, "Show version")
	OutputFlag  = flag.Bool("o", false, "Write to output file")
)

func Parse() {
	flag.Parse()
}

func PrintDefaults() {
	flag.PrintDefaults()
}
