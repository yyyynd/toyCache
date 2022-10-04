package cache

import (
	"container/list"
)

type LRUCache struct {
	maxEntries  int64                    //max cache entries bound
	curEntries  int64                    //had cached entries
	store       *list.List               //store entries
	reflectForm map[string]*list.Element //

}

type elementLRU struct {
	key string
	value interface{}
}

func NewLRUCache(maxEntries int64) Cache{
	return &LRUCache{
		maxEntries:  maxEntries,
		curEntries:  0,
		store:       list.New(),
		reflectForm: make(map[string]*list.Element),
	}
}

func (c *LRUCache) Add(key string, value interface{}){
	if c.curEntries == c.maxEntries {
		c.RemoveOldest()
	}
	ele := c.store.PushFront(&elementLRU{key : key,
									value : value})
	c.reflectForm[key] = ele
	c.curEntries += 1
}

func (c *LRUCache) Get(key string)(value interface{}, ok bool)  {
	if ele, ok := c.reflectForm[key]; !ok {
		return nil, false
	}else {
		c.store.MoveToFront(ele)
		return ele.Value.(*elementLRU).value, true
	}
}

func (c *LRUCache) Remove(key string) (removeKey string, removeValue interface{}){
	if c.curEntries == 0 {
		return
	}
	ele, ok := c.reflectForm[key]
	if !ok {
		return
	}
	entry := ele.Value.(*elementLRU)

	c.store.Remove(ele)
	delete(c.reflectForm, key)
	c.curEntries -= 1
	return entry.key, entry.value
}

func (c *LRUCache) RemoveOldest() (removeKey string, removeValue interface{}){
	//check element number
	if c.curEntries == 0 {
		return
	}
	ele := c.store.Back()
	entry := ele.Value.(*elementLRU)

	c.store.Remove(ele)
	delete(c.reflectForm, entry.key)
	c.curEntries -= 1
	return entry.key, entry.value
}

func (c *LRUCache) CacheSize() int64 {
	return c.curEntries
}

func (c *LRUCache) MaxCacheNum() int64 {
	return c.maxEntries
}