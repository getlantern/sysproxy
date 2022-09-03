package sysproxy

import (
	_ "embed"
	"fmt"
	"net"
	"os/exec"
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
	if len(sysproxy) == 0 {
		return fmt.Errorf("unable to find binary: %v", sysproxy)
	}
	be, err = byteexec.New(sysproxy, path)
	if err != nil {
		return fmt.Errorf("unable to extract helper tool: %v", err)
	}
	return ensureElevatedOnDarwin(be, prompt, iconFullPath)
}

// On tells OS to configure proxy through `addr` as host:port. It always returns
// a function that can be used to clear the system proxy setting. If the current
// process terminates before the clear function is called, the system proxy
// setting will be cleared anyway.
func On(addr string) (func() error, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse address %v: %v", addr, err)
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.To4() == nil {
		host = "[" + host + "]"
	}

	mu.Lock()
	defer mu.Unlock()
	if be == nil {
		return nil, fmt.Errorf("call EnsureHelperToolPresent() first")
	}

	cmd := be.Command("on", host, port)
	onErr := run(cmd)
	off, offErr := waitAndCleanup(host, port)
	if offErr != nil {
		log.Errorf("Unable to prepare waitAndCleanup job: %v", offErr)
	}
	if onErr != nil {
		return off, onErr
	}
	verifyErr := verify(addr)
	return off, verifyErr
}

// Off immediately unsets the proxy at addr as the system proxy.
func Off(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("unable to parse address %v: %v", addr, err)
	}

	mu.Lock()
	defer mu.Unlock()
	if be == nil {
		return fmt.Errorf("call EnsureHelperToolPresent() first")
	}

	cmd := be.Command("off", host, port)
	if err := run(cmd); err != nil {
		return err
	}
	return verify("")
}

type resultType struct {
	out []byte
	err error
}

func waitAndCleanup(host string, port string) (func() error, error) {
	cmd := be.Command("wait-and-cleanup", host, port)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	// Set up the command to run as a detached process
	detach(cmd)
	resultCh := make(chan *resultType)
	go func() {
		out, err := cmd.CombinedOutput()
		resultCh <- &resultType{
			out: out,
			err: err,
		}
	}()
	return func() error {
		stdin.Close()
		result := <-resultCh
		if result.err != nil {
			return fmt.Errorf("unable to finish %v: %s\n%s", cmd.Path, result.err, string(result.out))
		}
		return verify("")
	}, nil
}

func run(cmd *exec.Cmd) error {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to execute %v: %s\n%s", cmd.Path, err, string(out))
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
	actual := string(out)
	log.Debugf("Command %v output %v", cmd.Path, actual)
	if !allEquals(expected, actual) {
		return fmt.Errorf("unexpected output: expect '%s', got '%s'", expected, actual)
	}
	return nil
}

func allEquals(expected string, actual string) bool {
	if (expected == "") != (strings.TrimSpace(actual) == "") { // XOR
		return false
	}
	lines := strings.Split(actual, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" && trimmed != expected {
			return false
		}
	}
	return true
}
