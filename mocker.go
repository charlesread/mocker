// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {

	command := "/bin/sh"

	env := []string{
		"PS1=shell > ",
	}

	cloneFlags := syscall.CLONE_NEWNS |
		syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWIPC |
		syscall.CLONE_NEWPID |
		syscall.CLONE_NEWNET |
		syscall.CLONE_NEWUSER

	sysProcAttr := &syscall.SysProcAttr{
		Cloneflags: uintptr(cloneFlags),
	}

	cmd := exec.Command(command)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	cmd.SysProcAttr = sysProcAttr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error encountered running command %q: %s\n", command, err.Error())
		os.Exit(1)
	}

}
