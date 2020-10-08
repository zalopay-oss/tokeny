package totp

import (
	"encoding/base32"
	"encoding/binary"
	"strings"
	"time"

	"github.com/zalopay-oss/tokeny/pkg/hotp"
)

const (
	tokenLength        = 6
	bufferLength       = 8
	otpTTL       int64 = 30
)

type generator struct {
	secret []byte
}

func NewGenerator(secret string) (*generator, error) {
	s, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return nil, err
	}
	return &generator{s}, nil
}

func (g *generator) Generate() Token {
	now := time.Now().Unix()
	quotient := now / otpTTL
	remainder := now % otpTTL

	data := make([]byte, bufferLength)
	binary.BigEndian.PutUint64(data, uint64(quotient))

	return Token{
		Value:      hotp.Generate(g.secret, data, tokenLength),
		TimeoutSec: otpTTL - remainder,
	}
}
