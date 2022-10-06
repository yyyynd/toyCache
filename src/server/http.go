package server

import (
	"net/http"
	"strings"
	"toyCache/src/cache"
)

type GroupPool struct {
	self string		//this node add
	basePath string		//represent api
	groups map[string]*GroupNode
}

const defaultBasePath = "/groupCache/"

func NewGroupPool(self string, basePath string) *GroupPool {
	newPool := GroupPool{
		self: self,
		basePath: defaultBasePath,
	}
	if basePath != "" {
		newPool.basePath = basePath
	}
	return &newPool
}


func (g *GroupPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, g.basePath) {
		http.Error(w, "Wrong request", http.StatusNotFound)
		return
	}
	parts := strings.SplitN(r.URL.Path[len(g.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Wrong number of parameters", http.StatusBadRequest)
		return
	}
	//groupName := parts[0]
	key := parts[1]

	//w.Write([]byte(fmt.Sprintf("groupName: %s, key: %s", groupName, key)))
	//temporary
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(g.Get(key).(string)))
}

func (g *GroupPool) Get(key string) (value interface{}){
	//first check local cache
	if value, ok := cache.Get(key); ok {
		return value
	}
	//second check weather get from peers
	return nil
}






