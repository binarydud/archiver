package archive

import (
	"crypto/sha256"
	"encoding/hex"
)

func genSha(data []byte) (string, error) {
	h256 := sha256.New()
	h256.Write([]byte(data))
	shaSum := h256.Sum(nil)
	// sha256base64 := base64.StdEncoding.EncodeToString(shaSum[:])
	shaString := hex.EncodeToString(shaSum)

	return shaString, nil
}
