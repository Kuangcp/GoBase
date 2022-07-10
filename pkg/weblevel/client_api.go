package weblevel

import "github.com/kuangcp/gobase/pkg/ctool"

type Client interface {
	Del(key string) error
	Get(key string) (ctool.ResultVO[string], error)
	Set(key, val string)
	Sets(kv map[string]string)
	Stats()
	PrefixSearch(prefix string) map[string]string
}
