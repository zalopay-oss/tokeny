package hotp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	result := Generate([]byte("my-very-key"), []byte("my-very-counter"), 6)
	assert.Equal(t, "789672", result)
}

