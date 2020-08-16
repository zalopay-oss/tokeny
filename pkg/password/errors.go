package password

import "github.com/pkg/errors"

var (
	ErrPasswordsMismatch = errors.New("passwords do not match")
	ErrNotRegistered = errors.New("have not registered yet")
	ErrWrongPassword = errors.New("wrong password")
)
