package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

func main() {
	result, err := totp("")
	if err != nil {
		log.Panic(err)
	}
	log.Print(result)
}

func totp(secret string) (string, error) {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	quotient := now / 30
	//remainder := now % 30

	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(quotient))

	result := hotp(key, data, 6)

	return result, nil
}

func hotp(key []byte, counter []byte, otpLength int) string {
	// Generate 20-byte SHA1 hash
	hash := hmacSHA1(key, counter)
	// Get last half-byte to use as truncate offset
	offset := hash[19] & 0xf
	// Truncate 4 bytes
	truncatedHash := hash[offset : offset+4]
	// Remove first bit
	truncatedHash[0] = truncatedHash[0] & 0x7f
	// Convert to number
	hNum := binary.BigEndian.Uint32(truncatedHash)
	// Get last n-digit as OTP
	otp := hNum % (uint32)(math.Pow10(otpLength))
	// Padding with 0s
	result := padding0(fmt.Sprint(otp), otpLength)

	return result
}

func hmacSHA1(key []byte, counter []byte) []byte {
	h := hmac.New(sha1.New, key)
	h.Write(counter)
	return h.Sum(nil)
}

func padding0(otp string, expectedLength int) string {
	fmtTemplate := fmt.Sprintf("%%0%ds", expectedLength)
	return fmt.Sprintf(fmtTemplate, otp)
}
