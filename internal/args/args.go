package args

import (
	"flag"
	"fmt"
	"os"

	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/setup"
)

func Parse() bool {
	helpFlag := flag.Bool("h", false, "Display help")
	initFlag := flag.Bool("i", false, "Init project")
	versionFlag := flag.Bool("v", false, "Show version")
	flag.Parse()
	if *helpFlag {
		fmt.Println("HarzMind Code")
		fmt.Println("Usage: hzmind [flags] [commands]")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		return true
	}
	if *versionFlag {
		fmt.Fprintf(os.Stdout, "v%s\n", internal.VERSION_DATE)
		return true
	}
	if *initFlag {
		err := setup.SetupWorkingDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Fprintf(os.Stdout, "Project initiated")
	}
	return false
}
