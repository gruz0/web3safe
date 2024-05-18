package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gruz0/web3safe/internal/utils"
)

var (
	ErrCommandAlreadyRegistered = errors.New("command already registered")
	ErrUnknownCommand           = errors.New("unknown command")
)

type CommandError struct {
	name     string
	message  string
	exitCode int
}

func NewCommandError(name, message string, exitCode int) *CommandError {
	return &CommandError{
		name:     name,
		message:  message,
		exitCode: exitCode,
	}
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("%s: %v", e.name, e.message)
}

func (e *CommandError) ExitCode() int {
	return e.exitCode
}

type CommandHandler interface {
	ParseArgs(args []string) error
	Run() (bool, error)
	Name() string
	PrintHelp()
}

type Runner struct {
	commands map[string]CommandHandler
}

func NewRunner() *Runner {
	return &Runner{
		commands: make(map[string]CommandHandler),
	}
}

func (r *Runner) Register(command CommandHandler) error {
	if _, ok := r.commands[command.Name()]; ok {
		return fmt.Errorf("%w: %s", ErrCommandAlreadyRegistered, command.Name())
	}

	r.commands[strings.ToLower(command.Name())] = command

	return nil
}

func (r *Runner) Help() {
	fmt.Fprintf(os.Stdout, "Usage: %s subcommand arguments\n", utils.AppName)

	for name, cmd := range r.commands {
		fmt.Fprintf(os.Stdout, "\nSubcommand: %s\n", name)
		cmd.PrintHelp()
	}
}

func (r *Runner) Run(args []string) (bool, error) {
	commandName := strings.ToLower(args[0])

	cmd, ok := r.commands[commandName]
	if !ok {
		return false, fmt.Errorf("%w: %s", ErrUnknownCommand, commandName)
	}

	if err := cmd.ParseArgs(args[1:]); err != nil {
		return false, fmt.Errorf("%w", err)
	}

	success, err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	return success, nil
}
