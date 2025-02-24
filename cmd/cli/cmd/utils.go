package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func checkRootPermissions(msgs ...string) {
	msg := "You need root privileges to run this program\nPlease run with sudo\n"
	if len(msgs) != 0 {
		msg = msgs[0] + "\n"
	}
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
		fmt.Printf(msg + "\n\tsudo hybr init\n\n")
		os.Exit(1)
	}
}
