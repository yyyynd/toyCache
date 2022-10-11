package server

import (
	"hash/crc32"
)

type GroupNode struct {
	nodeName string
	address string	//add + port
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