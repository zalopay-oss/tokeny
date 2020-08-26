package crypto

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeHMACSHA1(t *testing.T) {
	result := ComputeHMACSHA1([]byte("my-something"), []byte("thequickbrownfoxjumpsoverthelazydog"))
	assert.Equal(t, "e5028c8a840936e8565629b795fb1c1e0bc53b0d", hex.EncodeToString(result))
}

