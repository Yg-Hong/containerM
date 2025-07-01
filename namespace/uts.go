package namespace

import (
	"containerM/utils"
	"syscall"
)

func SetHostname(name string) {
	utils.Must(syscall.Sethostname([]byte(name)), "Set hostname")
}
