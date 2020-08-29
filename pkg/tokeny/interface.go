package tokeny

import "github.com/zalopay-oss/tokeny/pkg/totp"

type Repository interface {
	Add(alias string, secret string) error
	Generate(alias string) (totp.Token, error)
	Delete(alias string) error
	List() ([]string, error)
	LastValidEntry() (string, error)
}
