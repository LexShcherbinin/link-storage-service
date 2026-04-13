package cache

type Cache interface {
	Get(string) (string, error)
	Set(string, string) error
	Delete(string) error
}
