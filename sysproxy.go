package sysproxy

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/getlantern/byteexec"
	"github.com/getlantern/golog"
)

var (
	log = golog.LoggerFor("sysproxy")

	mu sync.Mutex
	be *byteexec.Exec
)

// EnsureHelperToolPresent checks if helper tool exists and extracts it if not.
// On Mac OS, it also checks and set the file's owner to root:wheel and the setuid bit,
// it will request user to input password through a dialog to gain the rights to do so.
// path: absolute or relative path of the file to be checked and generated if
// not exists. Note - relative paths are resolved relative to the system-
// specific folder for aplication resources.
// prompt: the message to be shown on the dialog.
// iconPath: the full path of the icon to be shown on the dialog.
func EnsureHelperToolPresent(path string, prompt string, iconFullPath string) (err error) {
	mu.Lock()
	defer mu.Unlock()
	assetName := "sysproxy"
	// Load different binaries for 32bit and 64bit Windows respectively.
	if runtime.GOOS == "windows" {
		suffix := ".exe"
		assetName = assetName + suffix
	}
	sysproxyBytes, err := Asset(assetName)
	if err != nil {
		return fmt.Errorf("Unable to access sysproxy asset: %v", err)
	}
	be, err = byteexec.New(sysproxyBytes, path)
	if err != nil {
		return fmt.Errorf("Unable to extract helper tool: %v", err)
	}
	return ensureElevatedOnDarwin(be, prompt, iconFullPath)
}

// On tells OS to configure proxy through `addr` as host:port
func On(addr string) (err error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("Unable to parse address %v: %v", addr, err)
	}

	mu.Lock()
	defer mu.Unlock()
	if be == nil {
		return fmt.Errorf("call EnsureHelperToolPresent() first")
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = be.Command(addr)
	} else {
		cmd = be.Command("on", host, port)
	}

	if err := run(cmd); err != nil {
		return err
	}
	return verify(addr)
}

// Off sets proxy mode back to direct/none
func Off(addr string) (err error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("Unable to parse address %v: %v", addr, err)
	}

	mu.Lock()
	defer mu.Unlock()
	if be == nil {
		return fmt.Errorf("call EnsureHelperToolPresent() first")
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = be.Command()
	} else {
		cmd = be.Command("off", host, port)
	}
	if err := run(cmd); err != nil {
		return err
	}
	return verify(addr)
}

func run(cmd *exec.Cmd) error {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Unable to execute %v: %s\n%s", cmd.Path, err, string(out))
	}
	log.Debugf("Command %v output %v", cmd.Path, string(out))
	return nil
}

func verify(expected string) error {
	cmd := be.Command("show")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	str := string(out)
	log.Debugf("Command %v output %v", cmd.Path, str)
	if expected == "" && str != "" {
		return fmt.Errorf("Unexpected output %s", str)
	}
	lines := strings.Split(str, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" && trimmed != expected {
			return fmt.Errorf("Unexpected output %s", l)
		}
	}
	return nil
}
