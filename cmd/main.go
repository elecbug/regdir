package main

import (
	"fmt"
	"os"

	"github.com/elecbug/regdir/internal/cli"
	"github.com/elecbug/regdir/internal/cli/parameter"
)

func main() {
	// Parse command-line arguments and execute the corresponding command
	params := parameter.Parse(os.Args...)
	result, success := cli.Run(params)

	if success {
		fmt.Println(result)
	} else {
		fmt.Println("Error:", result)
		os.Exit(1)
	}
}
