package totp

type Token struct {
	Value      string
	TimeoutSec int64
}
