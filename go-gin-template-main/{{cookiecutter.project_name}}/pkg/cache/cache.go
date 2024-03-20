package cache

type Cache interface {
	Set(key, value interface{}, ttl int64, infinity bool) error
	Get(key interface{}) (interface{}, error)
}
