package sysproxy

import (
	_ "embed"
	"os/exec"
	"syscall"
)

//go:embed binaries/windows/sysproxy_386.exe
var sysproxy []byte

func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
