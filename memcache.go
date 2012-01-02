package memcache

type MemMap struct {
	items map[string]*MemMapItem
}

type MemMapItem struct {
	Raw []byte
	Gzip []byte
	Expire []byte
	Key string
}

func (m *MemMap)InitMap() {
	m.items = make(map[string] *MemMapItem, 100); 
}

func (m *MemMap) Add(item *MemMapItem) {
	m.items[item.Key] = item;
}

func (m *MemMap) GetByKey(key string) *MemMapItem {
	return m.items[key];
}