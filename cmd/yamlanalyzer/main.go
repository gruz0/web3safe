package main

import (
	"fmt"
	"os"

	"github.com/gruz0/web3safe/internal/config"
	"github.com/gruz0/web3safe/internal/yamlanalyzer"
)

func main() {
	flags := yamlanalyzer.ParseFlags()

	cfg := loadConfig(flags.ConfigFilePath)

	yamlAnalyzer := yamlanalyzer.NewYamlAnalyzer(flags.PathToScan, cfg)

	if err := yamlAnalyzer.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run yamlAnalyzer: %v\n", err)

		os.Exit(1)
	}

	yamlAnalyzerReport := yamlAnalyzer.Report()

	if len(yamlAnalyzerReport) == 0 {
		fmt.Fprintf(os.Stdout, "Nothing found. Great!\n")

		os.Exit(0)
	}

	for _, message := range yamlAnalyzerReport {
		fmt.Fprintln(os.Stderr, message)
	}

	os.Exit(1)
}

func loadConfig(configFilePath string) config.Config {
	if configFilePath == "" {
		fmt.Fprintf(os.Stdout, "No config file provided. We will use the default configuration.\n\n")

		return config.GetDefaultConfig()
	}

	fmt.Fprintf(os.Stdout, "Loading configuration file: %s\n", configFilePath)

	loadedConfig, loadConfigErr := config.LoadConfig(configFilePath)
	if loadConfigErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to load config: %v\n", loadConfigErr)
		os.Exit(1)
	}

	return loadedConfig
}
