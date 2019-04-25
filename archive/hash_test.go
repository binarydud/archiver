package archive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenFileShas(t *testing.T) {
	assert := assert.New(t)

	data := []byte("test")
	sha, err := genSha(data)
	assert.NoError(err, "should not error out")
	assert.Equal("9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", sha, "should be equal")
}
