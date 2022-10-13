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
	selfId      int    //this node hash key
	selfAddress string //this node address
	selfName    string //this node name
	basePath    string //represent api
	nodesId     []int
	groups      map[int]*GroupNode
	hit         int
	miss        int
	load		int
	LoadFunc    LoadFun
}

var LoadFunc Loader
var mu sync.Mutex
const defaultBasePath = "/groupCache/"

func NewGroupPool(selfAddress string, selfName string, basePath string) *GroupPool {
	newPool := GroupPool{
		selfAddress: selfAddress,
		selfName:    selfName,
		selfId:      consistentHash(selfAddress + selfName),
		basePath:    defaultBasePath,
		groups:      make(map[int]*GroupNode),
		nodesId:     make([]int, 0),
	}
	newPool.AddNode(selfName, selfAddress, newPool.selfId)
	if basePath != "" {
		newPool.basePath = basePath
	}
	LoadFunc = nil
	return &newPool
}

func (p *GroupPool) AddNode(nodeName string, address string, nodeId int) error {
	node := NewGroupNode(nodeName, address)
	//id := consistentHash(nodeName + address)
	if _, ok := p.groups[nodeId]; ok {
		return errors.New("node had existed")
	}
	p.groups[nodeId] = node
	p.nodesId = append(p.nodesId, nodeId)
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
	value := p.Get(key)
	if value == nil{
		w.WriteHeader(http.StatusBadRequest)
	}else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(value.(string)))
	}
}

func (p *GroupPool) Get(key string) (value interface{}){
	//first check local cache
	if value, ok := cache.Get(key); ok {
		p.hit += 1
		return value
	}
	p.miss += 1
	//second check weather get from peers
	//if this request need current node handle
	if peerId,err := p.peerPeek(key); err != nil{
		//log
	}else if peerId == p.selfId{
		p.load += 1
		//log
		//try load data by user specified method
		value, _ = LoadFunc.Load(key)
		return value
	}else {
		peer := p.groups[peerId]
		res, err := http.Get(peer.address + p.basePath + peer.nodeName + "/" +key)
		if err != nil {
			return nil
		}else {
			value, err = ioutil.ReadAll(res.Body)
			if err != nil {
				//log
				return nil
			}else {
				builder := strings.Builder{}
				for _, c := range value.([]uint8){
					builder.WriteString(string(c))
				}
				return builder.String()
			}
		}//else
	}//else
	return nil
}

func (p *GroupPool) peerPeek(key string) (peerId int, err error){
	id := consistentHash(key)
	for _, peerId = range p.nodesId {
		if id <= peerId {
			 return peerId, nil
		}
	}
	return p.nodesId[0], nil
}






