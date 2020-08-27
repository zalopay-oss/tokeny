package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPadding0(t *testing.T) {
	tests := []struct {
		input string
		desiredLength int
		expected string
	}{
		{
			input: "123456",
			desiredLength: 6,
			expected: "123456",
		},
		{
			input: "",
			desiredLength: 6,
			expected: "000000",
		},
		{
			input: "1",
			desiredLength: 6,
			expected: "000001",
		},
		{
			input: "321",
			desiredLength: 6,
			expected: "000321",
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := Padding0(test.input, test.desiredLength)
			assert.Equal(t, test.expected, result)
		})
	}
}

