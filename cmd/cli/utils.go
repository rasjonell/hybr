package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func checkRootPermissions() {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	i, err := strconv.Atoi(string(output[:len(output)-1]))

	if err != nil {
		panic(err)
	}

	if i != 0 {
		fmt.Printf(("You need root privileges to run this program\nPlease run with sudo\n"))
		os.Exit(1)
	}
}
