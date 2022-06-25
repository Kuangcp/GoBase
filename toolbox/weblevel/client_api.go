package weblevel

type Client interface {
	Del(key string) error
	Get(key string) (string, error)
	Set(key, val string)
	Sets(kv map[string]string)
}
