package sysproxy

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	assert.True(t, in("127.0.0.1:8888", "127.0.0.1:8888\n127.0.0.1:8888"))
	assert.True(t, in("127.0.0.1:8888", "127.0.0.1:8888\n127.0.0.1:8888\n"))
	assert.False(t, in("127.0.0.1:8888", "127.0.0.1:8888\n127.0.0.1:8887"))
	assert.True(t, in("", "\n\n"))
	assert.True(t, in("", "\r\n"))
	assert.False(t, in("", "127.0.0.1:8888"))
	assert.False(t, in("127.0.0.1:8888", ""))
}

func TestGetOutput(t *testing.T) {
	path := path.Join(os.TempDir(), "sysproxy")
	err := EnsureHelperToolPresent(path, "For test purpose", "")
	assert.NoError(t, err, "should install helper tool")
	err = On("localhost:8888")
	assert.NoError(t, err, "should set system proxy on")
	err = Off("localhost:8889")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unexpected output", "should fail to set system proxy off with correct address")
	err = Off("localhost:8888")
	assert.NoError(t, err, "should set system proxy off with correct address")
}
