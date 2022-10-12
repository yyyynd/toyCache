package server

import (
	"hash/crc32"
)

type GroupNode struct {
	nodeName string
	address string	//add + port
	hit int64
	miss int64
}

type hash func(data []byte) uint32
var hashFn = crc32.ChecksumIEEE

func SetHash(fn hash) {
	if fn == nil{
		return
	}
	hashFn = fn
}

func NewGroupNode(nodeName string, address string)*GroupNode {
	return &GroupNode{
				nodeName: nodeName,
				address: address}
}

//return a 0~2^32-1 int
func consistentHash(key string) int{
	return int(hashFn([]byte(key)))
}

func (n *GroupNode) HitStatistics() int64{
	return n.hit
}

func (n *GroupNode) MissStatistics() int64 {
	return n.miss
}

func (n *GroupNode) HitCount() {
	n.hit += 1
}

func (n *GroupNode) MissCount() {
	n.miss += 1
}