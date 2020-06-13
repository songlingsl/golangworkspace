package cmap

import "sync"

type cmap struct {
	Innermap map[string]string
	Mutex    *sync.RWMutex
}

func NewCmap() *cmap {
	return &cmap{Innermap: make(map[string]string), Mutex: new(sync.RWMutex)}
}

func (cmap *cmap) Put(key string, v string) {
	cmap.Mutex.Lock()
	defer cmap.Mutex.Unlock()
	cmap.Innermap[key] = v

}

func (cmap *cmap) Get(key string) string {
	cmap.Mutex.RLock()
	defer cmap.Mutex.RUnlock()
	return cmap.Innermap[key]
}
