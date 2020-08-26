package hotp

import (
	"encoding/binary"
	"fmt"
	"github.com/ltpquang/tokeny/pkg/crypto"
	"github.com/ltpquang/tokeny/pkg/utils"
	"math"
)

func Generate(key []byte, counter []byte, otpLength int) string {
	hash := crypto.ComputeHMACSHA1(key, counter)
	offset := hash[19] & 0xf
	truncatedHash := hash[offset : offset+4]
	truncatedHash[0] = truncatedHash[0] & 0x7f
	hNum := binary.BigEndian.Uint32(truncatedHash)
	otp := hNum % (uint32)(math.Pow10(otpLength))
	result := utils.Padding0(fmt.Sprint(otp), otpLength)
	return result
}
