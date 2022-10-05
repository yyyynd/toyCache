package cache

import (
	"errors"
	"runtime"
	"sync"
)

type Cache interface {
	add(key string, value interface{})
	get(key string)(value interface{}, ok bool)
	remove(key string) (removeKey string, removeValue interface{})
	currentEntries() int64
	maxCacheNum() int64
	removeOldest() (removeKey string, removeValue interface{})
	modifyMaxEntries(newMaxNum int64)
}

var cacheModeMap map[string]func(int64) Cache
var controller Cache
var mu sync.Mutex
const defaultMaxCacheNum = 5000

func init() {
	cacheModeMap  = make(map[string]func(int64) Cache)
	cacheModeMap["LRU"] = NewLRUCache
	cacheModeMap["LFU"] = NewLFUCache
	controller = NewLRUCache(defaultMaxCacheNum)
}

//SetCacheMode change cache mode
func SetCacheMode(modeName string, maxCacheNum int64) error {
	mu.Lock()
	defer mu.Unlock()

	if modeName == "" {
		return errors.New("mode name is nil")
	}else if maxCacheNum <= 0{
		return errors.New("maxCacheNum must > 0")
	} else if _, ok := cacheModeMap[modeName]; !ok {
		return errors.New("this mode not existed")
	}
	controller = cacheDataTransfer(cacheModeMap[modeName](maxCacheNum))
	runtime.GC()
	return nil
}

func Add(key string, value interface{}) {
	mu.Lock()
	defer mu.Unlock()
	controller.add(key, value)
}

func Get(key string) (value interface{}, ok bool){
	mu.Lock()
	defer mu.Unlock()
	value, ok = controller.get(key)
	return
}

func Remove(key string) (removeKey string, removeValue interface{}) {
	mu.Lock()
	defer mu.Unlock()
	removeKey, removeValue = controller.remove(key)
	return
}

func MaxCacheNum() int64 {
	mu.Lock()
	defer mu.Unlock()
	return controller.maxCacheNum()
}

func CurrentCacheNum() int64 {
	mu.Lock()
	defer mu.Unlock()
	return controller.currentEntries()
}

func RemoveOldest() (removeKey string, removeValue interface{}) {
	mu.Lock()
	defer mu.Unlock()
	removeKey, removeValue = controller.removeOldest()
	return
}

func ModifyMaxEntries(newMaxNum int64)  {
	mu.Lock()
	defer mu.Unlock()
	controller.modifyMaxEntries(newMaxNum)
}

func RemoveAll()  {
	mu.Lock()
	defer mu.Unlock()
	for i := controller.currentEntries(); i > 0; i-- {
		controller.removeOldest()
	}
}

//func AddCacheMode(modeName string, f func(int64) Cache) error{
//	if modeName == ""{
//		return errors.New("mode name is nil")
//	}else if f == nil {
//		return errors.New("f is nil")
//	}else if _, ok := cacheModeMap[modeName]; ok {
//		return errors.New("this mode existed")
//	}
//
//	cacheModeMap[modeName] = f
//	return nil
//}
//
//func ModifyExistCacheMode() error {
//	return nil
//}

func cacheDataTransfer(newCache Cache) Cache{
	//first check weather need delete
	removeNum := controller.currentEntries() - newCache.maxCacheNum()
	for ; removeNum > 0; removeNum-- {
		controller.removeOldest()
	}
	for i := controller.currentEntries(); i > 0; i-- {
		key, value := controller.removeOldest()
		newCache.add(key, value)
	}
	return newCache
}








