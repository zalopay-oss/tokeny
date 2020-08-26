package keyvalue

type Store interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	GetAllWithPrefixed(keyPrefix string) ([]KeyValue, error)
}
