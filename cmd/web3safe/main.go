package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/gruz0/web3safe/internal/app"
	"github.com/gruz0/web3safe/internal/utils"
)

var (
	goVersion = "unknown" //nolint:gochecknoglobals

	// Populated by goreleaser during build.
	version = "unknown"
	commit  = "?" //nolint:gochecknoglobals
	date    = ""  //nolint:gochecknoglobals
)

func main() {
	if len(os.Args) == 1 {
		printVersion()
		os.Exit(utils.ExitSuccess)
	}

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

func printVersion() {
	buildInfo, available := debug.ReadBuildInfo()
	if !available {
		return
	}

	var revision string

	for _, setting := range buildInfo.Settings {
		// The `vcs.xxx` information is only available with `go build`.
		// This information is not available with `go install` or `go run`.
		switch setting.Key {
		case "vcs.time":
			date = setting.Value
		case "vcs.revision":
			revision = setting.Value
		}
	}

	if revision == "" {
		revision = "unknown"
	}

	if date == "" {
		date = "(unknown)"
	}

	goVersion = buildInfo.GoVersion
	version = buildInfo.Main.Version
	commit = revision

	fmt.Fprintf(
		os.Stdout,
		"web3safe has version %s built with %s from %s on %s\n",
		version,
		goVersion,
		commit,
		date,
	)
}
