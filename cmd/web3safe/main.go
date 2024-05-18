package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gruz0/web3safe/internal/app"
	"github.com/gruz0/web3safe/internal/utils"
)

func main() {
	success, err := app.Run(os.Args)
	if err != nil {
		var appError *app.Error

		if errors.As(err, &appError) {
			fmt.Fprintln(os.Stderr, appError.Error())
			os.Exit(appError.ExitCode())
		}

		fmt.Fprintf(os.Stderr, "Unhandled error: %v", err)
		os.Exit(utils.ExitError)
	}

	if success {
		os.Exit(utils.ExitSuccess)
	}

	os.Exit(utils.ExitError)
}
