package server

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

var strByte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var strByteLen = len(strByte)

func RandString(length int) string {
	bytes := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		bytes[i] = strByte[r.Intn(strByteLen)]
	}
	return *(*string)(unsafe.Pointer(&bytes))
}

//need improve this test unit's policy
func TestPeerCacheGet(t *testing.T) {
	groupPool1 := NewGroupPool("http://localhost:4000", "node1","/groupCache/")
	groupPool2 := NewGroupPool("http://localhost:4001", "node2","/groupCache/")
	_ = groupPool1.AddNode("node2","http://localhost:4001", groupPool2.selfId )


	valueSet := make([]string, 2048)
	for i := 0; i < 2048; i++ {
		valueSet[i] = RandString(10)
		c1.Add(valueSet[i], valueSet[i])
		c2.Add(valueSet[i], valueSet[i])
	}
	//you need create a new goroutine to run listen
	go http.ListenAndServe(":4001", groupPool2)

	basicURL := "http://localhost:4000/groupCache/node1/"
	for i := 0; i < len(valueSet); i++ {
		req := httptest.NewRequest("GET", basicURL+valueSet[i], nil)
		w := httptest.NewRecorder()
		groupPool1.ServeHTTP(w, req)
		bytes, _ := ioutil.ReadAll(w.Result().Body)
		if valueSet[i] != *(*string)(unsafe.Pointer(&bytes)) {
			t.Fatalf("get fail, at: %d", i)
		}
	}

	t.Logf("hit: %d" ,groupPool1.groups[groupPool1.selfId].hit)


}
