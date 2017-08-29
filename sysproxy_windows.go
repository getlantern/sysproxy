// +build windows

package sysproxy

import (
	"os/exec"
	"syscall"
)

const DETACHED_PROCESS = 0x00000008
const CREATE_NEW_PROCESS_GROUP = 0x00000200

func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: CREATE_NEW_PROCESS_GROUP,
	}
}
