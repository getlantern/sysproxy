//go:build !darwin

package sysproxy

import (
	"github.com/getlantern/byteexec"
)

func ensureElevatedOnDarwin(be *byteexec.Exec, prompt string, iconFullPath string) (err error) {
	return nil
}
