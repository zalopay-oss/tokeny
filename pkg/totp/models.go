package totp

type Token struct {
	value      string
	timeoutSec int64
}
