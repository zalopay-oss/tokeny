package session

type Manager interface {
	IsSessionValid(sessionKey string) (bool, error)
	NewSession(sessionKey string) error
}
