// Package cli will initialize cli for config.
package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	cmd *cobra.Command
)

func init() {
	cmd = SetConfigCmds()
}

// CliMain will take the workload of executing/starting the cli, when the command is passed to it.
func CliMain() {
	err := Execute(os.Args[1:])
	if err != nil {
		cm.NeuronSaysItsError(err.Error())
		os.Exit(1)
	}
}

// Execute will actually execute the cli by taking the arguments passed to cli.
func Execute(args []string) error {

	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	if err != nil {
		return err
	}
	return nil
}
