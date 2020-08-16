package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
)

func ComputeHMACSHA1(key []byte, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	h.Write(data)
	return h.Sum(nil)
}
