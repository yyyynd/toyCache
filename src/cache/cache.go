package cache

type Cache interface {
	Add(key string, value interface{})
	Get(key string)(value interface{}, ok bool)
	Remove(key string)
	CacheSize() int64
}

var _ map[string]string

var CacheController = NewLRUCache(defaultMaxBytes) // default use LRU

//mb, kb, b
const defaultMaxBytes = 5000

//func SpecifyCacheSize(size int64) (int64, error){
//
//}








