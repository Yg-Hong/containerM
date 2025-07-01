package namespace

import (
	"containerM/utils"
	"fmt"
	"syscall"
)

func MountProc() {
	utils.Must(syscall.Mount("proc", "/proc", "proc", 0, ""), "mount /proc")
}

func MountSys() {
	utils.Must(syscall.Mount("sysfs", "/sys", "sysfs", 0, ""), "mount /sys")
}

func MountDev() {
	utils.Must(syscall.Mount("tmpfs", "/dev", "tmpfs", 0, ""), "mount /dev")

	// 디바이스 생성 함수
	createDev := func(path string, mode uint32, major, minor int) {
		dev := (major << 8) | minor
		err := syscall.Mknod(path, mode, dev)
		if err != nil {
			fmt.Printf("Failed to create %s: %v\n", path, err)
		}
		syscall.Chmod(path, 0666)
	}

	createDev("/dev/null", syscall.S_IFCHR|0666, 1, 3)
	createDev("/dev/zero", syscall.S_IFCHR|0666, 1, 5)
	createDev("/dev/random", syscall.S_IFCHR|0666, 1, 8)
	createDev("/dev/urandom", syscall.S_IFCHR|0666, 1, 9)
	createDev("/dev/tty", syscall.S_IFCHR|0666, 5, 0)
	createDev("/dev/console", syscall.S_IFCHR|0600, 5, 1)
}

func UnmountProc() {
	utils.Must(syscall.Unmount("/proc", 0))
}

func UnmountSys() {
	utils.Must(syscall.Unmount("/sys", 0))
}

func UnmountDev() {
	utils.Must(syscall.Unmount("/dev", 0))
}
