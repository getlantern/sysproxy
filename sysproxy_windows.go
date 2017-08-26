// +build windows

package sysproxy

import (
	"os/exec"
)

func detach(cmd *exec.Cmd) {
	// on Windows, we don't have to do anything special to detach process
}
