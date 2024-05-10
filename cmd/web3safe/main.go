package main

import (
	"fmt"
	"os"

	"github.com/gruz0/web3safe/internal/config"
	"github.com/gruz0/web3safe/internal/flags"
)

func main() {
	flags := flags.ParseFlags()

	if flags.GenerateConfig {
		generateConfig()
	}
}

func generateConfig() {
	newConfigFilePath := config.GetDefaultConfigFilePath()

	if err := config.GenerateConfig(newConfigFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating config: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "New configuration file generated at %s\n", newConfigFilePath)
	os.Exit(0)
}
