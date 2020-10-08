package tokeny

import "github.com/pkg/errors"

var (
	ErrNoEntryFound       = errors.New("no entry found")
	ErrEntryExistedBefore = errors.New("entry existed before")
)
