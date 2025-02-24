package cmd

import (
	"fmt"

	"hybr/cmd/cli/initiate"
	"hybr/internal/nginx"
	"hybr/internal/services"
	"os"

	"github.com/spf13/cobra"
)

type initFlagsType struct {
	email  string
	domain string

	forceNoSSL          bool
	forceResetTemplates bool
}

var initFlags initFlagsType

func init() {
	generateCmd.Flags().BoolVarP(
		&initFlags.forceResetTemplates,
		"forceDefaults", "f", false,
		"Reset default templates",
	)
	generateCmd.Flags().BoolVar(
		&initFlags.forceNoSSL,
		"no-ssl", false,
		"Don't use SSL",
	)
	generateCmd.Flags().StringVarP(
		&initFlags.email,
		"email", "e", "",
		"Specify Your Email for SSL certificate generation",
	)
	generateCmd.Flags().StringVarP(
		&initFlags.email,
		"domain", "d", "",
		"Specify Base Domain Name",
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
	services.InitRegistry(initFlags.forceResetTemplates)

	initiate.InitCLI(initFlags.email, initFlags.domain, initFlags.forceNoSSL)
	if _, err := initiate.NewProgram().Run(); err != nil {
		os.Exit(1)
	}
	model := initiate.GetModel()

	if !model.Done {
		fmt.Println("Service Installation Cancelled.")
		fmt.Println("Nothing to do.")
		os.Exit(0)
	}

	if err := nginx.Init(model.FinalBaseConfig, initFlags.forceResetTemplates, initFlags.forceNoSSL); err != nil {
		panic(err)
	}

	if err := services.InstallServices(model.GetFinalServices(), model.FinalBaseConfig); err != nil {
		panic(err)
	}
}
