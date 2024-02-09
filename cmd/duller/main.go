package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/anjolaoluwaakindipe/duller/internal/gateway"
)

// Runner is an interface for various commandline subsets for duller
type Runner interface {
	Name() string
	Init([]string) error
	Run() error
	UsageInfo()
}

// root takes in os command line arguments and invokes a subset command
// corresponding to subset that was called
func root(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Must provide a sub command")
	}

	subCommand := args[0]

	subCmds := []Runner{
		discovery.NewDiscCommand(),
		gateway.NewGateCommand(),
	}

	for _, subCmd := range subCmds {
		if subCmd.Name() == subCommand {
			if err := subCmd.Init(args[1:]); err != nil {
				return err
			}

			if err := subCmd.Run(); err != nil {
				return err
			}

			return nil
		}
	}

	for _, subCmd := range subCmds {
		subCmd.UsageInfo()
	}

	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		return nil
	}

	return fmt.Errorf("Unknown subcommand")
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		slog.Error(fmt.Sprintf("%v", err))
	}
}
