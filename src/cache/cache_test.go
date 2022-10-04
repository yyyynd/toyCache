package cache

import (
	"reflect"
	"testing"
)


var entrySet = map[string]string{"key1":"a","key2":"b", "key3":"c"}
func init()  {
	CacheController.Add("key1","a")
	CacheController.Add("key2","b")
	CacheController.Add("key3","c")
}

func TestSetCacheMode(t *testing.T) {
	t.Logf(reflect.Indirect(reflect.ValueOf(CacheController)).Type().Name())
	if err := SetCacheMode("LFU", defaultMaxCacheNum); err != nil{
		t.Fatalf("Set Fail")
	}
	
}

