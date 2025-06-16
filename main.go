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
		Cloneflags: syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

// Exec in child process
func child() {
	fmt.Printf("Running child process PID %d\n", syscall.Getpid())

	must(syscall.Sethostname([]byte("containerM")), "set hostname")

	// Configure new root
	newRoot := "/tmp/ubuntufs"
	putOld := newRoot + "/.pivot_root"

	// Mount point isolation
	must(syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""), "make mount private")

	// bind mount new root to itself (required by pivot_root)
	must(syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""), "bind mount newRoot")

	// make .pivot_root directory
	must(os.MkdirAll(putOld, 0700), "create putOld")

	// pivot_root(new_root, put_old)
	must(syscall.PivotRoot(newRoot, putOld), "pivot_root")

	// change working directory to new root
	must(os.Chdir("/"), "chdir to /")

	// unmount old root
	must(syscall.Unmount("/.pivot_root", syscall.MNT_DETACH), "unmount old root")
	must(os.RemoveAll("/.pivot_root"), "remove old root")

	// mount /proc, /sys, /dev
	must(syscall.Mount("proc", "/proc", "proc", 0, ""), "mount /proc")
	must(syscall.Mount("sysfs", "/sys", "sysfs", 0, ""), "mount /sys")
	must(syscall.Mount("tmpfs", "/dev", "tmpfs", 0, ""), "mount /dev")

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
