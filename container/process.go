package container

import (
	"containerM/fs"
	"containerM/namespace"
	"containerM/utils"

	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Run(command []string) {
	fmt.Printf("Parent: Running %v as PID %d\n", command, os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, command...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// isolation namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET,
	}

	// setting Environment variable
	cmd.Env = append(os.Environ(), "SETUP_NET=1")

	utils.Must(cmd.Start())

	namespace.SetupHostNet(cmd)
	// wait for child process to finish
	utils.Must(cmd.Wait())

	defer namespace.CleanupHostNet(cmd)
}

func Child(command []string) {
	fmt.Printf("Running child process PID %d\n", syscall.Getpid())

	// set hostname
	namespace.SetHostname("containerM")

	// set file system pivot root
	fs.SetupPivotRoot("/tmp/ubuntufs")

	// mount & cleanup /dev, /proc, /sys
	namespace.MountDev()
	namespace.MountProc()
	namespace.MountSys()
	defer namespace.UnmountDev()
	defer namespace.UnmountProc()
	defer namespace.UnmountSys()

	// network configuration
	if os.Getenv("SETUP_NET") == "1" {
		pid := fmt.Sprintf("%d", syscall.Getpid())
		cmd := exec.Command("/setup_container_net.sh", pid)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		utils.Must(cmd.Run(), "run setup_net.sh")
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	utils.Must(cmd.Run())
}
