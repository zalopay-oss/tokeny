package password

type Manager interface {
	IsRegistered() (bool, error)
	Register(pwd string, rePwd string) error
	Login(pwd string) error
}
