package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestTOTP(t *testing.T) {
	tests := []struct {
		key string
	}{
		{"YUKWZZN7YEC5FTSF"},
	}
	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			otp, err := totp(test.key)
			assert.NoError(t, err)
			log.Print(otp)
		})
	}
}
