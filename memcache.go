package memcache

import (
	"fmt"
	"sync"
	"http"
)

var RepoLock sync.RWMutex;

type MemMap struct {
	Items map[string]*MemMapItem
}

type MemMapItem struct {
	Raw []byte
	Gzip []byte
	Expire []byte
	Key string
	Head http.Header
}

func (m *MemMap)InitMap() {
	m.Items = make(map[string] *MemMapItem, 100);
}

func (m *MemMap) Add(v *MemMapItem) {
	RepoLock.Lock();
	m.Items[v.Key] = v;
	RepoLock.Unlock();
	fmt.Print("+")
}

func (m *MemMap) GetByKey(key string) (*MemMapItem) {
	RepoLock.RLock();
	defer RepoLock.RUnlock();
	if mapRet, ok := m.Items[key]; ok {
		fmt.Print(".")
		return mapRet;
	}
 
	fmt.Print("-")
	return nil;
}

