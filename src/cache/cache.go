package cache

import (
	"errors"
	"runtime"
)

type Cache interface {
	Add(key string, value interface{})
	Get(key string)(value interface{}, ok bool)
	Remove(key string) (removeKey string, removeValue interface{})
	CacheSize() int64
	MaxCacheNum() int64
	RemoveOldest() (removeKey string, removeValue interface{})
}

var cacheModeMap map[string]func(int64) Cache
var CacheController Cache
const defaultMaxCacheNum = 5000

func init() {
	cacheModeMap  = make(map[string]func(int64) Cache)
	cacheModeMap["LRU"] = NewLRUCache
	cacheModeMap["LFU"] = NewLFUCache
	CacheController = NewLRUCache(defaultMaxCacheNum)
}

//func SpecifyCacheSize(size int64) (int64, error){
//
//}

//SetCacheMode change cache mode
func SetCacheMode(modeName string, maxCacheNum int64) error {
	if modeName == "" {
		return errors.New("mode name is nil")
	}else if maxCacheNum <= 0{
		return errors.New("maxCacheNum must > 0")
	} else if _, ok := cacheModeMap[modeName]; !ok {
		return errors.New("this mode not existed")
	}
	CacheController = cacheDataTransfer(cacheModeMap[modeName](maxCacheNum))
	runtime.GC()
	return nil
}

func AddCacheMode(modeName string, f func(int64) Cache) error{
	if modeName == ""{
		return errors.New("mode name is nil")
	}else if f == nil {
		return errors.New("f is nil")
	}else if _, ok := cacheModeMap[modeName]; ok {
		return errors.New("this mode existed")
	}

	cacheModeMap[modeName] = f
	return nil
}

func ModifyExistCacheMode() error {
	return nil
}

func cacheDataTransfer(newCache Cache) Cache{
	//first check weather need delete
	removeNum := CacheController.CacheSize() - newCache.MaxCacheNum()
	for ; removeNum > 0; removeNum-- {
		CacheController.RemoveOldest()
	}
	for i := CacheController.CacheSize(); i > 0; i-- {
		key, value := CacheController.RemoveOldest()
		newCache.Add(key, value)
	}

	return newCache
}








