package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"toyCache/src/cache"
)

type GroupPool struct {
	self     string //this node address
	selfId   int    //this node hash key
	selfName string //this node name
	basePath string //represent api
	nodesId  []int
	groups   map[int]*GroupNode
	mu sync.Mutex
}

const defaultBasePath = "/groupCache/"

func NewGroupPool(self string, selfName string, basePath string) *GroupPool {
	newPool := GroupPool{
		self:     self,
		selfName: selfName,
		selfId: consistentHash(self + selfName),
		basePath: defaultBasePath,
		groups:   make(map[int]*GroupNode),
		nodesId:  make([]int, 1024),
	}
	if basePath != "" {
		newPool.basePath = basePath
	}
	return &newPool
}

func (p *GroupPool) AddNode(nodeName string, address string) error {
	node := NewGroupNode(nodeName, address)
	id := consistentHash(nodeName + address)
	if _, ok := p.groups[id]; ok {
		return errors.New("node had existed")
	}
	p.groups[id] = node
	p.nodesId = append(p.nodesId, id)
	sort.Ints(p.nodesId)
	return nil
}


func (p *GroupPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		http.Error(w, "Wrong request", http.StatusNotFound)
		return
	}
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Wrong number of parameters", http.StatusBadRequest)
		return
	}
	//groupName := parts[0]
	key := parts[1]

	//w.Write([]byte(fmt.Sprintf("groupName: %s, key: %s", groupName, key)))
	//temporary
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(p.Get(key).(string)))
}

func (p *GroupPool) Get(key string) (value interface{}){
	//first check local cache
	if value, ok := cache.Get(key); ok {
		return value
	}
	//second check weather get from peers
	//if this request need current node handle
	//res, err := http.Get("http://localhost:4000/groupCache/node2/test2")
	//if err != nil {
	//	return nil
	//}else{
	//	value, err = ioutil.ReadAll(res.Body)
	//	if err != nil {
	//		return nil
	//	}else {
	//		return value
	//	}
	//}

	if peerId,err := p.peerPeek(key); err != nil{
		return nil
	}else if peerId == p.selfId{
		//try load data from database
	}else {
		peer := p.groups[peerId]
		res, err := http.Get(peer.address +"/"+ p.basePath+"/" + peer.nodeName + "/" +key)
		if err != nil {
			return nil
		}else {
			value, err = ioutil.ReadAll(res.Body)
			if err != nil {
				return nil
			}else {
				return value
			}
		}//else
	}//else
	return nil
}

func (p *GroupPool) peerPeek(key string) (peerId int, err error){
	id := consistentHash(key)
	for peerId = range p.nodesId {
		if id > peerId {
			 return peerId, nil
		}
	}
	return 0,errors.New("no peers")
}






