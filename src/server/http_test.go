package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"toyCache/src/cache"
	"unsafe"

	c1 "toyCache/src/cache"
	c2 "toyCache/src/cache"
)

// Attempt to read locally stored data
func TestLocalCacheGet(t *testing.T) {
	groupPool := NewGroupPool("","groupTest", "/groupCache/")
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

//need improve this test unit's policy
func TestPeerCacheGet(t *testing.T) {
	c1.Add("test1","test1")
	c2.Add("test2","test2")
	groupPool1 := NewGroupPool("http://localhost:4000", "node1","/groupCache/")
	groupPool2 := NewGroupPool("http://localhost:4001", "node2","/groupCache/")
	_ = groupPool1.AddNode("node2","http://localhost:4001" )

	//you need create a new goroutine to run listen
	go http.ListenAndServe(":4001", groupPool2)

	req := httptest.NewRequest("GET","http://localhost:4000/groupCache/node1/test2", nil)
	w := httptest.NewRecorder()
	groupPool1.ServeHTTP(w, req)
	bytes, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatalf("fail")
	}else {
		t.Logf("%s",bytes)
	}

}
