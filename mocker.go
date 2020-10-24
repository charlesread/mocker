// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"path/filepath"

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

func pivotRoot(newroot string) error {
	putold := filepath.Join(newroot, "/.pivot_root")

	// bind mount newroot to itself - this is a slight hack
	// needed to work around a pivot_root requirement
	if err := syscall.Mount(
		newroot,
		newroot,
		"",
		syscall.MS_BIND|syscall.MS_REC,
		"",
	); err != nil {
		return err
	}

	// create putold directory
	if err := os.MkdirAll(putold, 0700); err != nil {
		return err
	}

	// call pivot_root
	if err := syscall.PivotRoot(newroot, putold); err != nil {
		return err
	}

	// ensure current working directory is set to new root
	if err := os.Chdir("/"); err != nil {
		return err
	}

	// umount putold, which now lives at /.pivot_root
	putold = "/.pivot_root"
	if err := syscall.Unmount(
		putold,
		syscall.MNT_DETACH,
	); err != nil {
		return err
	}

	// remove putold
	if err := os.RemoveAll(putold); err != nil {
		return err
	}

	return nil
}

func nsInitialisation() {
	fmt.Printf("\n>> namespace setup code goes here <<\n\n")
	newrootPath := os.Args[1]

	if err := pivotRoot(newrootPath); err != nil {
		fmt.Printf("Error running pivot_root - %s\n", err)
		os.Exit(1)
	}
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

	cmd := reexec.Command("nsInitialisation", "./root")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = sysProcAttr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the reexec.Command - %s\n", err)
		os.Exit(1)
	}

}
