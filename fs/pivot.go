package fs

import (
	"containerM/utils"
	"os"
	"os/exec"
	"syscall"
)

func SetupPivotRoot(newRoot string) {
	putOld := newRoot + "/.pivot_root"

	// Mount point isolation
	utils.Must(syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""), "make mount private")

	// bind mount new root to itself (required by pivot_root)
	utils.Must(syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""), "bind mount newRoot")

	// make .pivot_root directory
	utils.Must(os.MkdirAll(putOld, 0700), "create putOld")

	// copy script file
	utils.Must(exec.Command("cp", "./scripts/setup_container_net.sh", newRoot+"/setup_container_net.sh").Run(), "copy setup_net.sh")

	// pivot_root(new_root, put_old)
	utils.Must(syscall.PivotRoot(newRoot, putOld), "pivot_root")

	// change working directory to new root
	utils.Must(os.Chdir("/"), "chdir to /")

	// unmount old root
	utils.Must(syscall.Unmount("/.pivot_root", syscall.MNT_DETACH), "unmount old root")
	utils.Must(os.RemoveAll("/.pivot_root"), "remove old root")
}
