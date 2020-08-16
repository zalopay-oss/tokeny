package totp

type Generator interface {
	Generate() Token
}
