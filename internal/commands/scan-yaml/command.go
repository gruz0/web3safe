package scanyaml

import (
	"flag"
	"fmt"
	"os"

	"github.com/gruz0/web3safe/internal/commands"
	"github.com/gruz0/web3safe/internal/config"
	"github.com/gruz0/web3safe/internal/utils"
	"github.com/gruz0/web3safe/internal/yamlanalyzer"
)

type Command struct {
	fs              *flag.FlagSet
	configFilePath  string
	directoryToScan string
	fileToScan      string
	recursive       bool
}

func New() *Command {
	command := &Command{
		fs:              flag.NewFlagSet("yaml", flag.ContinueOnError),
		configFilePath:  "",
		directoryToScan: "",
		fileToScan:      "",
		recursive:       false,
	}

	command.fs.StringVar(&command.configFilePath, "config", "", "Path to configuration file")
	command.fs.StringVar(&command.directoryToScan, "dir", "", "Path to the directory to scan")
	command.fs.StringVar(&command.fileToScan, "file", "", "Path to the file to scan")
	command.fs.BoolVar(&command.recursive, "recursive", false, "Scan directory recursively")

	return command
}

func (c *Command) Name() string {
	return c.fs.Name()
}

func (c *Command) PrintHelp() {
	c.fs.PrintDefaults()
}

func (c *Command) ParseArgs(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return commands.NewCommandError(
			c.Name(),
			fmt.Sprintf("failed to parse args: %v", err),
			utils.ExitInvalidArguments,
		)
	}

	return nil
}

func (c *Command) Run() (bool, error) {
	if c.directoryToScan == "" && c.fileToScan == "" {
		c.PrintHelp()

		return false, nil
	}

	cfg, err := c.loadConfig()
	if err != nil {
		return false, err
	}

	yamlAnalyzer := yamlanalyzer.NewYamlAnalyzer(cfg)

	if c.directoryToScan != "" {
		if err := yamlAnalyzer.ScanDirectory(c.directoryToScan, c.recursive); err != nil {
			return false, commands.NewCommandError(
				c.Name(),
				fmt.Sprintf("Unable to analyze yaml files: %v", err),
				utils.ExitError,
			)
		}
	}

	if c.fileToScan != "" {
		if err := yamlAnalyzer.ScanFile(c.fileToScan); err != nil {
			return false, commands.NewCommandError(
				c.Name(),
				fmt.Sprintf("Unable to analyze yaml file: %v", err),
				utils.ExitError,
			)
		}
	}

	yamlAnalyzerReport := yamlAnalyzer.Report()

	if len(yamlAnalyzerReport) == 0 {
		fmt.Fprintln(os.Stdout, "Nothing found. Great!")

		return true, nil
	}

	for _, message := range yamlAnalyzerReport {
		fmt.Fprintln(os.Stderr, message)
	}

	return false, nil
}

func (c *Command) loadConfig() (config.Config, error) {
	cfg := config.GetDefaultConfig()

	if c.configFilePath != "" {
		loadedConfig, err := config.Load(c.configFilePath)
		if err != nil {
			return cfg, commands.NewCommandError(
				c.Name(),
				err.Error(),
				utils.ExitError,
			)
		}

		cfg = loadedConfig
	}

	return cfg, nil
}
