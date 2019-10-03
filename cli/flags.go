package cli

import (
	"github.com/spf13/cobra"
)

var (
	jsonAuth gcloudAuth
)

// Registering all the flags to the command neuron itself.
func registerFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&jsonAuth.JSONPath, "json", "j", "", "path to gcp auth json file")
	cmd.PersistentFlags().StringVarP(&jsonAuth.k8clusterName, "cluster-name", "c", "", "Name of the cluster which needs to be connected to")
	cmd.PersistentFlags().StringVarP(&jsonAuth.version, "version", "v", "1", "version of the cluster")
}
