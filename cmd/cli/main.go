package main

import (
	"hybr/internal/nginx"
	"hybr/internal/services"
	"os"
)

func main() {
	// Checking root permissions here for now
	// In the future when the CLI does more than setup services
	// We will only ask for root permissions when needed
	checkRootPermissions()

	services.InitRegistry(flags.forceResetTemplates)

	InitCLI()
	if _, err := NewProgram().Run(); err != nil {
		os.Exit(1)
	}

	if len(model.finalServices) == 0 {
		os.Exit(0)
	}

	if err := nginx.Init(model.finalBaseConfig, flags.forceResetTemplates); err != nil {
		panic(err)
	}

	if err := services.InstallServices(model.finalServices, model.finalBaseConfig); err != nil {
		panic(err)
	}
}
