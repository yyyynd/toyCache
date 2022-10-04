package cache

import (
	"container/list"
	"testing"
)


var c =  &LFUCache{
		maxEntries:  defaultMaxCacheNum,
		curEntries:  0,
		store:       list.New(),
		reflectForm: make(map[string]*list.Element),
	}

func TestLFUCache_Add(t *testing.T) {
	c.Add("key1","a")
	c.Add("key2","b")
	c.Add("key3","c")
	if c.curEntries != 3 {
		t.Fatalf("Wrong number after adding")
	}
	entrySet := map[string]string{"key1":"a","key2":"b", "key3":"c"}
	for k, v := range entrySet{
		if value, ok := c.Get(k); !ok{
			t.Fatalf("Unsuccessful storage")
		}else if value.(string) != v{
			t.Fatalf("Stored data error")
		}
	}
}

func TestLFUCache_Get(t *testing.T) {
	_,_ = c.Get("key1")
	_,_ = c.Get("key3")
	_,_ = c.Get("key3")
	_,_ = c.Get("key1")
	_,_ = c.Get("key1")
	_,_ = c.Get("key2")
	//except order 1 3 2
	order := []string{"a","c","b"}
	for e , i := c.store.Front(), 0; e != nil && i < len(order);
			e ,i = e.Next(), i+1{
		if order[i] != e.Value.(*elementLFU).value.(string){
			t.Fatalf("Order error")
		}
		//t.Logf("Current element frequecy :%d\n",e.Value.(*elementLFU).frequency)
	}
}

func TestLFUCache_Remove(t *testing.T) {
	c.Remove("key3")
	if _, ok := c.Get("key3"); ok{
		t.Fatalf("Remove error")
	}
	t.Logf("cur cache size : %d\n", c.CacheSize())
}

func TestLFUCache_RemoveOldest(t *testing.T) {
	c.RemoveOldest()
	if _, ok := c.Get("key2"); ok{
		t.Fatalf("Remove oldest error")
	}
	t.Logf("cur cache size : %d\n", c.CacheSize())
}