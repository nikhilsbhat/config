package cli

import (
	"github.com/spf13/cobra"
)

var (
	cmds map[string]*cobra.Command
)

type confcmds struct {
	commands []*cobra.Command
}

// SetConfigCmds helps in gathering all the subcommands so that it can be used while registering it with main command.
func SetConfigCmds() *cobra.Command {
	cmd := getConfigCmds()
	return cmd
}

func getConfigCmds() *cobra.Command {

	var configCmd = &cobra.Command{
		Use:   "config [command]",
		Short: "command to deal with config activities",
		Long:  `This will help user to deal with gcloud and kube config activity.`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  cm.echoConfig,
	}
	configCmd.SetUsageTemplate(getUsageTemplate())

	var setCmd = &cobra.Command{
		Use:          "set [flags]",
		Short:        "command to set the config",
		Long:         `This will help user to set the configurations.`,
		Run:          setConfig,
		SilenceUsage: true,
	}

	// Creating "version" happens here.
	var versionCmd = &cobra.Command{
		Use:   "version [flags]",
		Short: "command to fetch the version of config installed",
		Long:  `This will help user to find what version of Config he/she installed in her machine.`,
		RunE:  versionConfig,
	}

	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(versionCmd)
	registerFlags(configCmd)
	return configCmd
}

func (cm *cliMeta) echoConfig(cmd *cobra.Command, args []string) error {
	cmd.Usage()
	return nil
}

// This function will return the custom template for usage function,
// only functions/methods inside this package can call this.

func getUsageTemplate() string {
	return `{{printf "\n"}}Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if gt (len .Aliases) 0}}{{printf "\n" }}
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}{{printf "\n" }}
Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{printf "\n"}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}{{printf "\n"}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}{{printf "\n"}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}{{printf "\n"}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}{{printf "\n"}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}"
{{printf "\n"}}`
}
