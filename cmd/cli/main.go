package main

import (
	"fmt"
	"os"
)

func main() {
	if _, err := NewProgram().Run(); err != nil {
		os.Exit(1)
	}

	for name, s := range model.selected {
		fmt.Printf("Vars for %s\n\n", name)
		for _, v := range s.Variables {
			fmt.Printf("%s = %s\n", v.Name, v.Input.Value())
		}
	}
}
