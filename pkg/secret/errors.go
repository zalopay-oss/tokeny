package secret

import "github.com/pkg/errors"

var (
	ErrNoSecret = errors.New("secret not found")
)
