package app

import (
	"errors"
	"fmt"

	"github.com/gruz0/web3safe/internal/commands"
	generateconfig "github.com/gruz0/web3safe/internal/commands/generate-config"
	scandotenv "github.com/gruz0/web3safe/internal/commands/scan-dotenv"
	scanuserenv "github.com/gruz0/web3safe/internal/commands/scan-shell-env"
	scanyaml "github.com/gruz0/web3safe/internal/commands/scan-yaml"
	"github.com/gruz0/web3safe/internal/utils"
)

const minimumArgsRequired = 2

func Run(args []string) (bool, error) {
	runner := commands.NewRunner()

	cmds := []commands.CommandHandler{
		generateconfig.New(),
		scanuserenv.New(),
		scandotenv.New(),
		scanyaml.New(),
	}

	for _, cmd := range cmds {
		if err := runner.Register(cmd); err != nil {
			return false, NewError(
				fmt.Sprintf("failed to register command: %v", err),
				utils.ExitError,
			)
		}
	}

	if len(args) < minimumArgsRequired {
		runner.Help()

		return false, nil
	}

	success, err := runner.Run(args[1:])
	if err != nil {
		var cmdErr *commands.CommandError

		if errors.As(err, &cmdErr) {
			return false, NewError(
				err.Error(),
				cmdErr.ExitCode(),
			)
		}

		return false, fmt.Errorf("Error: %w", err)
	}

	return success, nil
}
