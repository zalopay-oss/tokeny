package keyvalue

type Store interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}
