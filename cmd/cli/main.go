package main

import (
	"hybr/internal/services"
	"os"
)

func main() {
	services.InitRegistry()

	if _, err := NewProgram().Run(); err != nil {
		os.Exit(1)
	}
}
