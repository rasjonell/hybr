package main

import (
	"fmt"
	"hybr/internal/services"
	"os"
)

func main() {
	if _, err := NewProgram().Run(); err != nil {
		os.Exit(1)
	}

	if len(model.finalServices) == 0 {
		os.Exit(0)
	}

	if err := services.InstallServices(model.finalServices); err != nil {
		panic(err)
	}

	fmt.Println("Check out ~/.hybr/services")
}
