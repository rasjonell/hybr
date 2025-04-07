package cmd

import "github.com/spf13/cobra"

type serviceCmdFlag struct {
	service string
	remote  string
}

var serviceCmdFlags serviceCmdFlag

func init() {
	servicesCmd.Flags().StringVarP(
		&serviceCmdFlags.service,
		"service", "s", "",
		"Name of the service",
	)

	servicesCmd.PersistentFlags().StringVar(
		&serviceCmdFlags.remote,
		"host", "",
		"Hostname of a remote hybr node",
	)

	rootCmd.AddCommand(servicesCmd)
}

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Show services info",
	Long:  "Show installed services with their status",
	Run:   func(cmd *cobra.Command, args []string) { cmd.Usage() },
}
