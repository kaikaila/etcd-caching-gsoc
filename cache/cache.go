package cache
type Cache interface {
	Get(key String)(string, bool)
	Set(key string, value string)
}