package cache

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
}

func BenchmarkLRU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < defaultMaxCacheNum * 8; j++ {
			if j % 250 == 0 {
				rand.Seed(time.Now().UnixNano())
			}
			key := string(rune(rand.Intn(6000)))
			if _, ok := Get(key); !ok{
				Add(key, key)
			}
		}
	}
}

func BenchmarkLFU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < defaultMaxCacheNum * 8; j++ {
			if j % 250 == 0 {
				rand.Seed(time.Now().UnixNano())
			}
			key := string(rune(rand.Intn(6000)))
			if _, ok := Get(key); !ok{
				Add(key, key)
			}
		}
	}
}

func BenchmarkInOrder(b *testing.B) {
	b.Run("BenchmarkLRU", BenchmarkLFU )
	SetCacheMode("LFU", defaultMaxCacheNum)
	RemoveAll()
	b.Run("BenchmarkLFU", BenchmarkLFU)
}