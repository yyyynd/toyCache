package cache

import (
	"testing"
)

func newLRU() *LRUCache {
	return  NewLRUCache(defaultMaxBytes)
}


func TestLRU(t *testing.T) {
	//c := NewLRU(defaultMaxBytes).(*LRUCache)
	c := newLRU()
	c.Add("key1","a")
	c.Add("key2","b")
	c.Add("key3","c")
	//read cache size
	t.Logf("cur cache size : %d\n", c.CacheSize())
	//read test
	entrySet := map[string]string{"key1":"a","key2":"b", "key3":"c"}
	for k, v := range entrySet{
		if value, ok := c.Get(k); !ok{
			t.Fatalf("Cache miss")
		}else if value.(string) != v{
			t.Fatalf("Get error")
		}
	}
	t.Logf("Get test pass")

	//least recently used principle test
	_,_ = c.Get("key1")
	_,_ = c.Get("key3")
	_,_ = c.Get("key1")
	_,_ = c.Get("key2")
	//order need 2 1 3
	order := []string{"b","a","c"}
	for e , i := c.store.Front(), 0; e != nil && i < len(order);
			e ,i = e.Next(), i+1{
		if order[i] != e.Value.(elementLRU).value.(string){
			t.Fatalf("Order error")
		}
	}
	t.Logf("Order test pass")
	//remove test
	c.Remove("key1")
	if _, ok := c.Get("key1"); ok{
		t.Fatalf("Remove error")
	}
	t.Logf("Remove test pass")
	t.Logf("cur cache size : %d\n", c.CacheSize())
}