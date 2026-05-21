package cli

import (
	"fmt"

	"github.com/elecbug/regdir/internal/cli/parameter"
)

const VERSION = "0.1.0"
const AUTHOR = "elecbug"
const DESCRIPTION = "RegDir, A simple file name batch modifier"

// Run is the main entry point for the CLI application. It processes the command-line parameters and
// executes the appropriate actions based on the provided options.
func Run(params *parameter.Parameters) {
	if params.MainParam == "" {
		runHelp(params)
	} else {

	}
}

// runHelp displays the usage information for the CLI application.
func runHelp(params *parameter.Parameters) {
	result := `Usage: ` + params.MainParam + ` [OPTIONS]...`
	fmt.Println(result)
}
