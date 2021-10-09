package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker		  run <image> <cmd> <params>
// go run container.go /bin/bash

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("invalid command")
	}
}

func run() {
	fmt.Printf("running %v as PID %d", os.Args[2:], os.Getpid())

	// execute command as child process using isolated namespaces
	args := append([]string{"child"}, os.Args[2:]...)
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// new namespaces: UTS, PID, MNT
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {
	// set up new rootfs
	must(syscall.Chroot("/var/rootfs")) // standard linux fs needs to have been created already
	must(os.Chdir("/"))

	// set up separate dir for procs at the target proc with the type of proc and a data value of proc.
	must(syscall.Mount("proc", "proc", "proc", 0, "proc")) // mount proc to proc with type proc and data proc

	// run command
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
