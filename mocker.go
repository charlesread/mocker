package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error encountered running command: %s\n", err.Error())
		os.Exit(1)
	}

}
