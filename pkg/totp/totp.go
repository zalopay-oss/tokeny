package totp

import (
	"encoding/base32"
	"encoding/binary"
	"github.com/ltpquang/tokeny/pkg/hotp"
	"strings"
	"time"
)

func Generate(secret string) (string, error) {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	quotient := now / 30
	//remainder := now % 30

	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(quotient))

	result := hotp.Generate(key, data, 6)

	return result, nil
}
