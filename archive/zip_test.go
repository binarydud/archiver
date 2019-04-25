package archive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetZipArchiver(t *testing.T) {
	assert := assert.New(t)
	archiver := GetZipArchiver("testpath.zip")
	assert.Equal("testpath.zip", archiver.filepath)
}
