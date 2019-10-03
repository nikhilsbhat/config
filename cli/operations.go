package cli

import (
	"github.com/spf13/cobra"
)

func setConfig(cmd *cobra.Command, args []string) {

	if err := configSet(jsonAuth); err != nil {
		cm.NeuronSaysItsError(getStringOfMessage(err))
	}

}
