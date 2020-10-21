// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/reexec"
)

func init() {
	t := time.Now().UnixNano()
	fmt.Printf("%v - init() - PPID: %d\n", t, os.Getppid())
	fmt.Printf("%v - init() - PID: %d\n", t, os.Getpid())
	fmt.Printf("%v - init() - os.Args: %v\n", t, os.Args)
	reexec.Register("nsInitialisation", nsInitialisation)
	if i := reexec.Init(); i == true {
		fmt.Printf("%v - init() - reexec.Init() is %v, exiting\n", t, i)
		fmt.Printf("%v - ------\n", t)
		os.Exit(0)
	} else {
		fmt.Printf("%v - init() - reexec.Init() is %v, continuing\n", t, i)
		fmt.Printf("%v - ------\n", t)
	}
}

func nsInitialisation() {
	fmt.Printf("\n>> namespace setup code goes here <<\n\n")
	nsRun()
}

func nsRun() {

	env := []string{
		"PS1=`hostname` > ",
	}

	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = env

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}

func main() {

	cloneFlags := syscall.CLONE_NEWNS |
		syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWIPC |
		syscall.CLONE_NEWPID |
		syscall.CLONE_NEWNET |
		syscall.CLONE_NEWUSER

	uidMappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      os.Getuid(),
			Size:        1,
		},
	}

	gidMappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      os.Getgid(),
			Size:        1,
		},
	}

	sysProcAttr := &syscall.SysProcAttr{
		Cloneflags:  uintptr(cloneFlags),
		UidMappings: uidMappings,
		GidMappings: gidMappings,
	}

	cmd := reexec.Command("nsInitialisation")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = sysProcAttr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the reexec.Command - %s\n", err)
		os.Exit(1)
	}

}
