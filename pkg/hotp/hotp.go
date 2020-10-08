package hotp

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/zalopay-oss/tokeny/pkg/crypto"
	"github.com/zalopay-oss/tokeny/pkg/utils"
)

func Generate(key []byte, counter []byte, otpLength int) string {
	hash := crypto.ComputeHMACSHA1(key, counter)
	truncatedHash := truncate(hash)
	hNum := binary.BigEndian.Uint32(truncatedHash)
	otp := hNum % (uint32)(math.Pow10(otpLength))
	result := utils.Padding0(fmt.Sprint(otp), otpLength)
	return result
}

func truncate(hash []byte) []byte {
	var (
		last4BitFilter byte = 0xf
		firstBitFilter byte = 0x7f
	)
	offset := hash[19] & last4BitFilter
	truncatedHash := hash[offset : offset+4]
	truncatedHash[0] &= firstBitFilter
	return truncatedHash
}
