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


func NewLFUCache(maxEntries int64) Cache {
	return &LFUCache{
		maxEntries:  maxEntries,
		curEntries:  0,
		store:       list.New(),
		reflectForm: make(map[string]*list.Element),
	}
}

func (c *LFUCache) add(key string, value interface{}){
	if c.curEntries == c.maxEntries {
		c.removeOldest()
	}
	ele := c.store.PushBack(&elementLFU{key: key,
								frequency: 0,
								value: value})
	c.reflectForm[key] = ele
	c.curEntries += 1
}

func (c *LFUCache) get(key string)(value interface{}, ok bool)  {
	ele, ok := c.reflectForm[key]
	if !ok {
		return nil, false
	}
	ele.Value.(*elementLFU).frequency += 1
	insert := c.insertPosition(ele)
	if insert == nil {
		c.store.MoveToFront(ele)
	}else {
		c.store.MoveAfter(ele, insert)
	}

	return ele.Value.(*elementLFU).value, true
}

func (c *LFUCache) remove(key string) (removeKey string, removeValue interface{}){
	if c.curEntries == 0 {
		return
	}
	ele, ok := c.reflectForm[key]
	if !ok {
		return
	}
	entry := ele.Value.(*elementLFU)

	c.store.Remove(ele)
	delete(c.reflectForm, key)
	c.curEntries -= 1
	return entry.key, entry.value
}

func (c *LFUCache) currentEntries() int64 {
	return c.curEntries
}

//insertPosition Find insert position
func (c *LFUCache) insertPosition(targetEle *list.Element) *list.Element{
	if c.curEntries == 0 {
		return nil
	}
	targetFre := targetEle.Value.(*elementLFU).frequency
	for ele := targetEle.Prev(); ele != nil; ele = ele.Prev(){
		if ele.Value.(*elementLFU).frequency > targetFre{
			return ele
		}
	}
	//if all ele.frequency <= targetEle.frequency or none other ele
	return nil
}

func (c *LFUCache) removeOldest() (removeKey string, removeValue interface{}){
	if c.curEntries == 0 {
		return
	}
	ele := c.store.Back()
	entry := ele.Value.(*elementLFU)
	delete(c.reflectForm, entry.key)
	c.store.Remove(ele)
	c.curEntries -= 1
	return entry.key, entry.value
}

func (c *LFUCache) maxCacheNum() int64 {
	return c.maxEntries
}

func (c *LFUCache) modifyMaxEntries(newMaxNum int64) {
	diff := c.curEntries - newMaxNum
	for ; diff > 0; diff-- {
		c.removeOldest()
	}
	c.maxEntries = newMaxNum
}

