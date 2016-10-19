package sysproxy

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOutput(t *testing.T) {
	path := path.Join(os.TempDir(), "pac")
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
