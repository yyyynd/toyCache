package protection

import (
	"math/rand"
	"testing"
	"time"
)

// should uncomment part of the code
func TestSingle_Do(t *testing.T) {
	s := Single{set: map[string]*call{}}
	fn := func(key string)(interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return key, nil
	}

	key := []string{"a","b","c","d","e"}
	for i := 0; i < 1000*1000; i++{
		if i % 100 == 0 {
			rand.Seed(time.Now().UnixNano())
		}
		go func(key string) {
			v, _ := s.Do(key, fn)
			if v.(string) != key {
				t.Error("Get fail")
				return
			}
		}(key[rand.Intn(5)])
	}

	t.Logf("do count : %d", s.callCount)
}

