package totp

import (
	"encoding/base32"
	"encoding/binary"
	"github.com/zalopay-oss/tokeny/pkg/hotp"
	"strings"
	"time"
)

const (
	tokenLength = 6
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
	quotient := now / 30
	remainder := now % 30

	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(quotient))

	return Token{
		Value:      hotp.Generate(g.secret, data, tokenLength),
		TimeoutSec: 30 - remainder,
	}
}
