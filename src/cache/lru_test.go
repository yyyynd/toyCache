package cache

import (
	"container/list"
	"testing"
)

func newLRU() *LRUCache {
	return &LRUCache{
		maxEntries:  defaultMaxCacheNum,
		curEntries:  0,
		store:       list.New(),
		reflectForm: make(map[string]*list.Element),
	}
}


func TestLRU(t *testing.T) {
	//c := NewLRU(defaultMaxCacheNum).(*LRUCache)
	c := newLRU()
	c.add("key1","a")
	c.add("key2","b")
	c.add("key3","c")
	//read cache size
	t.Logf("cur cache size : %d\n", c.currentEntries())
	//read test
	entrySet := map[string]string{"key1":"a","key2":"b", "key3":"c"}
	for k, v := range entrySet{
		if value, ok := c.get(k); !ok{
			t.Fatalf("Cache miss")
		}else if value.(string) != v{
			t.Fatalf("get error")
		}
	}
	t.Logf("get test pass")

	//least recently used principle test
	_,_ = c.get("key1")
	_,_ = c.get("key3")
	_,_ = c.get("key1")
	_,_ = c.get("key2")
	//order need 2 1 3
	order := []string{"b","a","c"}
	for e , i := c.store.Front(), 0; e != nil && i < len(order);
			e ,i = e.Next(), i+1{
		if order[i] != e.Value.(*elementLRU).value.(string){
			t.Fatalf("Order error")
		}
	}
	t.Logf("Order test pass")
	//remove test
	c.remove("key1")
	if _, ok := c.get("key1"); ok{
		t.Fatalf("remove error")
	}
	t.Logf("remove test pass")
	t.Logf("cur cache size : %d\n", c.currentEntries())
}

//func TestStorageOrder(t *testing.T) {
//	keyList := list.New()
//	for i := 0; i < defaultMaxCacheNum; i++ {
//
//	}
//}