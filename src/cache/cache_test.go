package cache

import (
	"reflect"
	"testing"
)


var entrySet = map[string]string{"key1":"a","key2":"b", "key3":"c"}
func init()  {
	Add("key1","a")
	Add("key2","b")
	Add("key3","c")
}

func TestSetCacheMode(t *testing.T) {
	t.Logf(reflect.Indirect(reflect.ValueOf(controller)).Type().Name())
	if err := SetCacheMode("LFU", defaultMaxCacheNum); err != nil{
		t.Fatalf("Default mode lose")
	}
	if reflect.Indirect(reflect.ValueOf(controller)).Type().Name() != "LFUCache"{
		t.Fatalf("Set mode fail")
	}
	t.Logf(reflect.Indirect(reflect.ValueOf(controller)).Type().Name())
	for k, v := range entrySet{
		value, ok := Get(k)
		if !ok {
			t.Fatalf("Cache transfer fail")
		}else if value.(string) != v {
			t.Fatalf("Cache transfer fail")
		}
	}
}

func TestModifyMaxEntries(t *testing.T) {
	ModifyMaxEntries(2)
	if CurrentCacheNum() != 2 || len(controller.(*LFUCache).reflectForm) != 2 ||
		controller.(*LFUCache).store.Len() != 2{
		t.Fatalf("Modify fail")
	}
}

func TestRemoveAll(t *testing.T) {
	RemoveAll()
	if CurrentCacheNum() != 0 || len(controller.(*LFUCache).reflectForm) != 0 ||
		controller.(*LFUCache).store.Len() != 0{
		t.Fatalf("Remove all fail")
	}
}


