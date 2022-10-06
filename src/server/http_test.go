package server

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"toyCache/src/cache"
	"unsafe"
)

// Attempt to read locally stored data
func TestLocalCacheGet(t *testing.T) {
	groupPool := NewGroupPool("", "/groupCache/")
	cache.Add("keyTest","valueTest")

	req := httptest.NewRequest("GET",
		"http://test.com/groupCache/groupTest/keyTest",nil)
	w := httptest.NewRecorder()
	groupPool.ServeHTTP(w, req)

	bytes, _ := ioutil.ReadAll(w.Result().Body)
	if *(*string)(unsafe.Pointer(&bytes)) != "valueTest" {
		t.Fatal("Get fail")
	}
}
