package scanshellenv

import (
	"flag"
	"fmt"
	"os"

	"github.com/gruz0/web3safe/internal/commands"
	"github.com/gruz0/web3safe/internal/config"
	"github.com/gruz0/web3safe/internal/envanalyzer"
	"github.com/gruz0/web3safe/internal/utils"
)

type Command struct {
	fs             *flag.FlagSet
	configFilePath string
}

func New() *Command {
	command := &Command{
		fs:             flag.NewFlagSet("shellenv", flag.ContinueOnError),
		configFilePath: "",
	}

	command.fs.StringVar(&command.configFilePath, "config", "", "Path to the configuration file")

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
	cfg := config.GetDefaultConfig()

	if c.configFilePath != "" {
		loadedConfig, err := config.Load(c.configFilePath)
		if err != nil {
			return false, commands.NewCommandError(
				c.Name(),
				err.Error(),
				utils.ExitError,
			)
		}

		cfg = loadedConfig
	}

	envAnalyzer := envanalyzer.NewEnvAnalyzer(cfg)

	if err := envAnalyzer.Run(); err != nil {
		return false, commands.NewCommandError(
			c.Name(),
			fmt.Sprintf("Unable to analyze shell envs: %v", err),
			utils.ExitError,
		)
	}

	envAnalyzerReport := envAnalyzer.Report()

	if len(envAnalyzerReport) == 0 {
		fmt.Fprintln(os.Stdout, "Nothing found in ENV. Great!")

		return true, nil
	}

	for _, message := range envAnalyzerReport {
		fmt.Fprintln(os.Stderr, message)
	}

	return false, nil
}
