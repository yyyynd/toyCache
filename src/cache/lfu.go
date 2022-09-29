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
	ele := c.store.PushBack(&elementLFU{key: key,
								frequency: 1,
								value: value})
	c.reflectForm[key] = ele
	c.curEntries += 1
}

func (c *LFUCache) Get(key string)(value interface{}, ok bool)  {
	ele, ok := c.reflectForm[key]
	if !ok {
		return nil, false
	}
	ele.Value.(*elementLFU).frequency += 1
	insert := c.insertPosition(ele.Value.(*elementLFU).frequency)
	if insert == nil {
		c.store.MoveToFront(ele)
	}else {
		c.store.MoveAfter(ele, insert)
	}

	return ele.Value.(*elementLFU).value, true
}

func (c *LFUCache) Remove(key string) {
	if c.curEntries == 0 {
		return
	}
	ele, ok := c.reflectForm[key]
	if !ok {
		return
	}
	c.store.Remove(ele)
	delete(c.reflectForm, key)
	c.curEntries -= 1
}

func (c *LFUCache) CacheSize() int64 {
	return c.curEntries
}

//insertPosition Find insert position
func (c *LFUCache) insertPosition(target int) *list.Element{
	if c.curEntries == 0 {
		return nil
	}
	for ele := c.store.Front(); ele != nil; ele = ele.Next(){
		if ele.Value.(*elementLFU).frequency <= target{
			return ele.Prev()
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
	delete(c.reflectForm, ele.Value.(*elementLFU).key)
	c.store.Remove(ele)
	c.curEntries -= 1
}

