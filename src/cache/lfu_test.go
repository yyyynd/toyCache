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
	c.add("key1","a")
	c.add("key2","b")
	c.add("key3","c")
	if c.curEntries != 3 {
		t.Fatalf("Wrong number after adding")
	}
	entrySet := map[string]string{"key1":"a","key2":"b", "key3":"c"}
	for k, v := range entrySet{
		if value, ok := c.get(k); !ok{
			t.Fatalf("Unsuccessful storage")
		}else if value.(string) != v{
			t.Fatalf("Stored data error")
		}
	}
}

func TestLFUCache_Get(t *testing.T) {
	_,_ = c.get("key1")
	_,_ = c.get("key3")
	_,_ = c.get("key3")
	_,_ = c.get("key1")
	_,_ = c.get("key1")
	_,_ = c.get("key2")
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
	c.remove("key3")
	if _, ok := c.get("key3"); ok{
		t.Fatalf("remove error")
	}
	t.Logf("cur cache size : %d\n", c.currentEntries())
	t.Logf("cur map size :%d\n", len(c.reflectForm))
}

func TestLFUCache_RemoveOldest(t *testing.T) {
	c.removeOldest()
	if _, ok := c.get("key2"); ok{
		t.Fatalf("remove oldest error")
	}
	t.Logf("cur cache size : %d\n", c.currentEntries())
	t.Logf("cur map size :%d\n", len(c.reflectForm))
}