package protection

import "sync"

/**
To deal with the cache-breakdown problem, only one of the requests with the
same key will be processed for all requests coming in at the same time, and
the remaining requests will be processed by the `call` struct.
 */

type Single struct {
	set map[string]*call
	mu sync.Mutex
	callCount int // test use
}

type call struct {
	wg sync.WaitGroup
	value interface{}
	err error
}

func (s *Single) Do(key string, fn func (string)(interface{}, error))(interface{}, error){
	s.mu.Lock()
	if c, ok := s.set[key]; ok{
		s.mu.Unlock()
		c.wg.Wait()
		return c.value, c.err
	}
	s.callCount += 1	//test use
	c := new(call)
	s.set[key] = c
	c.wg.Add(1)
	s.mu.Unlock()

	c.value, c.err = fn(key)
	c.wg.Done()
	s.mu.Lock()
	delete(s.set, key)
	s.mu.Unlock()

	return c.value, c.err
}


