// nolint: gosec
package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"log"
)

func ComputeHMACSHA1(key []byte, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	_, err := h.Write(data)
	if err != nil {
		log.Fatalln(err)
	}
	return h.Sum(nil)
}
