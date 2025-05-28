package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: containerM run <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("unknown command")
	}
}

func run() {
	fmt.Printf("Parent: Running %v as PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// isolation namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID,
	}

	must(cmd.Run())
}

// Exec in child process
func child() {
	fmt.Printf("Running child process PID %d\n", syscall.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}

func must(err error, context ...string) {
	if err != nil {
		if len(context) > 0 {
			fmt.Fprintf(os.Stderr, "Error (%s): %v\n", context[0], err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
