package cache

import (
	"sync"
)

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat()
}
type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}
func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}

type inMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Stat
}

func (C *inMemoryCache) Set(k string, v []byte) error {
	C.mutex.Lock()
	defer C.mutex.Unlock()
	tmp, exist := C.c[k]
	if exist {
		C.del(k, tmp)
	}
	C.c[k] = v
	C.add(k, v)
	return nil
}
func (C *inMemoryCache) Get(k string) ([]byte, error) {
	C.mutex.RLock()
	defer C.mutex.RUnlock()
	return C.c[k], nil
}
func (C *inMemoryCache) Del(k string) error {
	C.mutex.Lock()
	defer C.mutex.Unlock()
	v, exist := C.c[k]
	if exist {
		delete(C.c, k)
		C.del(k, v)
	}
	return nil
}
func (C *inMemoryCache) GetStat() Stat {
	return C.Stat
}
func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		c:     make(map[string][]byte),
		mutex: sync.RWMutex{},
		Stat:  Stat{},
	}
}
