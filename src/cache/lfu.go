package cache

import "container/list"

type LFUCache struct {
	maxEntries int64
	curEntries int64
	store *list.List
	reflectForm map[string]*list.Element
}

type elementLFU struct {
	frequency int
	key string
	value interface{}
}

func NewLFUCache(maxEntries int64) *LFUCache {
	return &LFUCache{
		maxEntries:  maxEntries,
		curEntries:  0,
		store:       list.New(),
		reflectForm: make(map[string]*list.Element),
	}
}

func (c *LFUCache) Add(key string, value interface{}){
	if c.curEntries == c.maxEntries {
		c.removeOldest()
	}
	ele := c.store.PushBack(elementLFU{key: key,
								frequency: 1,
								value: value})
	c.reflectForm[key] = ele
	c.curEntries += 1
}

//insertPosition Find insert position
func (c *LFUCache) insertPosition(target int) *list.Element{
	if c.curEntries == 0 {
		return nil
	}

	for ele := c.store.Front(); ele != nil; ele = ele.Next(){
		if ele.Value.(elementLFU).frequency <= target{
			return ele
		}
	}
	//if all ele.frequency > target, return list tail
	return c.store.Back()
}

func (c *LFUCache) removeOldest() {
	if c.curEntries == 0 {
		return
	}
	ele := c.store.Back()
	delete(c.reflectForm, ele.Value.(elementLFU).key)
	c.store.Remove(ele)
	c.curEntries -= 1
}

