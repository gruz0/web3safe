package generateconfig

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gruz0/web3safe/internal/commands"
	"github.com/gruz0/web3safe/internal/config"
	"github.com/gruz0/web3safe/internal/utils"
	"gopkg.in/yaml.v2"
)

const (
	defaultConfigFilename = "web3safe.yaml"
	configDirPermission   = 0o755
	configPermission      = 0o600
)

type Command struct {
	fs             *flag.FlagSet
	printHelp      bool
	configFilePath string
	createConfig   bool
	printConfig    bool
	force          bool
}

func New() *Command {
	command := &Command{
		fs:             flag.NewFlagSet("config", flag.ContinueOnError),
		printHelp:      false,
		configFilePath: "",
		createConfig:   false,
		printConfig:    false,
		force:          false,
	}

	command.fs.StringVar(&command.configFilePath, "config", "", "Path to the configuration file")
	command.fs.BoolVar(&command.printHelp, "help", false, "Show help")
	command.fs.BoolVar(&command.createConfig, "create", false, "Create a config by its default location")
	command.fs.BoolVar(&command.printConfig, "print", false, "Print a config")
	command.fs.BoolVar(&command.force, "force", false, "Overwrite config")

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
	if c.createConfig {
		return c.handleCreateConfig(c.configFilePath)
	}

	if c.printConfig {
		return c.handlePrintConfig(c.configFilePath)
	}

	c.PrintHelp()

	return false, nil
}

func (c *Command) handleCreateConfig(configFilePath string) (bool, error) {
	if configFilePath == "" {
		configFilePath = filepath.Join(utils.GetDefaultConfigDirectory(), defaultConfigFilename)
	}

	if utils.IsFileExist(configFilePath) && !c.force {
		return false, commands.NewCommandError(
			c.Name(),
			fmt.Sprintf("file %s already exists. Please use -force flag to override the config", configFilePath),
			utils.ExitError,
		)
	}

	yamlData, err := yaml.Marshal(config.GetDefaultConfig())
	if err != nil {
		return false, commands.NewCommandError(
			c.Name(),
			fmt.Sprintf("failed to encode config: %v", err),
			utils.ExitError,
		)
	}

	configDir := filepath.Dir(configFilePath)

	if err := os.MkdirAll(configDir, configDirPermission); err != nil {
		return false, commands.NewCommandError(
			c.Name(),
			fmt.Sprintf("failed to create a config directory %s: %v", configDir, err),
			utils.ExitError,
		)
	}

	err = os.WriteFile(configFilePath, yamlData, configPermission)
	if err != nil {
		return false, commands.NewCommandError(
			c.Name(),
			fmt.Sprintf("failed to write file %s: %v", configFilePath, err),
			utils.ExitError,
		)
	}

	fmt.Fprintf(os.Stdout, "New configuration file saved to %s\n", configFilePath)

	return true, nil
}

func (c *Command) handlePrintConfig(configFilePath string) (bool, error) {
	cfg := config.GetDefaultConfig()

	if configFilePath != "" {
		loadedConfig, err := config.Load(configFilePath)
		if err != nil {
			return false, commands.NewCommandError(
				c.Name(),
				err.Error(),
				utils.ExitError,
			)
		}

		cfg = loadedConfig
	}

	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		return false, commands.NewCommandError(
			c.Name(),
			fmt.Sprintf("failed to encode config: %v", err),
			utils.ExitError,
		)
	}

	_, err = os.Stdout.Write(yamlData)
	if err != nil {
		return false, commands.NewCommandError(
			c.Name(),
			fmt.Sprintf("failed to write config: %v", err),
			utils.ExitError,
		)
	}

	return true, nil
}
