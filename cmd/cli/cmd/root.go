package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hybr [command]",
	Short: "hybr is self-hosted service management platform",
	Long: `Hybr - A Fast and Simple Self-Hosted Service Management Platform built with
       love by rasjonell in Go.
       Complete documentation is available at https://docs.hybr.dev`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
