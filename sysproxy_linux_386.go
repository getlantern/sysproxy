package sysproxy

import (
	_ "embed"
	"os/exec"
	"syscall"

	"github.com/getlantern/byteexec"
)

//go:embed binaries/linux_386/sysproxy
var sysproxy []byte

func ensureElevatedOnDarwin(be *byteexec.Exec, prompt string, iconFullPath string) (err error) {
	return nil
}

func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}
