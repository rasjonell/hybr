package main

import (
	"fmt"
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

	if !model.done {
		fmt.Println("Service Installation Cancelled.")
		fmt.Println("Nothing to do.")
		os.Exit(0)
	}

	if err := nginx.Init(model.finalBaseConfig, flags.forceResetTemplates, flags.forceNoSSL); err != nil {
		panic(err)
	}

	if err := services.InstallServices(model.getFinalServices(), model.finalBaseConfig); err != nil {
		panic(err)
	}
}
