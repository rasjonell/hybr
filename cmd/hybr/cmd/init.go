package cmd

import (
	"fmt"

	"github.com/rasjonell/hybr/cmd/hybr/initiate"
	"github.com/rasjonell/hybr/internal/services"
	"github.com/rasjonell/hybr/internal/tailscale"
	"os"

	"github.com/spf13/cobra"
)

type initFlagsType struct {
	authKey             string
	forceResetTemplates bool
}

var initFlags initFlagsType

func init() {
	generateCmd.Flags().BoolVarP(
		&initFlags.forceResetTemplates,
		"forceDefaults", "f", false,
		"Reset default templates",
	)

	generateCmd.Flags().StringVarP(
		&initFlags.authKey,
		"ts-auth", "a", "",
		"Tailscale AUTH_KEY",
	)

	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "init",
	Short: "initate hybr project",
	Long:  "Initates a hybr project with your selection of services and configuration.",
	Run:   runInit,
}

func runInit(cmd *cobra.Command, _ []string) {
	checkRootPermissions("root privileges are required for initiating a hybr project.")
	services.InitRegistry(initFlags.forceResetTemplates, initFlags.authKey)

	initiate.InitCLI()
	if _, err := initiate.NewProgram().Run(); err != nil {
		os.Exit(1)
	}
	model := initiate.GetModel()

	if !model.Done {
		fmt.Println("Service Installation Cancelled.")
		fmt.Println("Nothing to do.")
		os.Exit(0)
	}

	if err := tailscale.Start(initFlags.authKey); err != nil {
		panic(err)
	}

	if err := services.InstallServices(model.GetFinalServices()); err != nil {
		panic(err)
	}
}
