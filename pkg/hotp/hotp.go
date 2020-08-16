package hotp

import (
	"encoding/binary"
	"fmt"
	"github.com/ltpquang/tokeny/pkg/crypto"
	"github.com/ltpquang/tokeny/pkg/utils"
	"math"
)

func Generate(key []byte, counter []byte, otpLength int) string {
	// Generate 20-byte SHA1 hash
	hash := crypto.ComputeHMACSHA1(key, counter)
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
	result := utils.Padding0(fmt.Sprint(otp), otpLength)

	return result
}
