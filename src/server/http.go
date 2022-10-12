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

func NewGroupPool(selfAddre string, selfName string, basePath string) *GroupPool {
	newPool := GroupPool{
		self:     selfAddre,
		selfName: selfName,
		selfId:   consistentHash(selfAddre + selfName),
		basePath: defaultBasePath,
		groups:   make(map[int]*GroupNode),
		nodesId:  make([]int, 0),
	}
	newPool.AddNode(selfName, selfAddre, newPool.selfId)
	if basePath != "" {
		newPool.basePath = basePath
	}
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(p.Get(key).(string)))
}

func (p *GroupPool) Get(key string) (value interface{}){
	//first check local cache
	//if value, ok := cache.Get(key); ok {
	//	//p.groups[p.selfId].HitCount()
	//	return value
	//}
	//p.groups[p.selfId].MissCount()
	//second check weather get from peers
	//if this request need current node handle
	if peerId,err := p.peerPeek(key); err != nil{

	}else if peerId == p.selfId{
		//try load data from database
		value, _ := cache.Get(key)
		p.groups[p.selfId].HitCount()
		return value
	}else {
		p.groups[p.selfId].MissCount()
		peer := p.groups[peerId]
		res, err := http.Get(peer.address + p.basePath + peer.nodeName + "/" +key)
		if err != nil {
			return nil
		}else {
			value, err = ioutil.ReadAll(res.Body)
			if err != nil {
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






