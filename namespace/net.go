package namespace

import (
	"containerM/utils"
	"fmt"
	"os/exec"
)

func SetupHostNet(cmd *exec.Cmd) {
	// network setup: run host-side network setup script with child PID
	utils.Must(exec.Command("bash", "./scripts/setup_host_net.sh", fmt.Sprintf("%d", cmd.Process.Pid)).Run(), "setup host network configure")
}

func CleanupHostNet(cmd *exec.Cmd) {
	utils.Must(exec.Command("bash", "./scripts/cleanup_host_net.sh", fmt.Sprintf("%d", cmd.Process.Pid)).Run(), "cleanup host network configure")
}
