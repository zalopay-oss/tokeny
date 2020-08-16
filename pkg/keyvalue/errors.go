package keyvalue

import "github.com/pkg/errors"

var (
	ErrNoRecord = errors.New("record not found")
)
