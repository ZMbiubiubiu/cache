package cache

import (
	"log"
	"sync"

	"cache/helper"
)

// LocalGetter 用户自己实现各自加载数据功能，用于查找缓存失败时从本地加载
type LocalGetter interface {
	LocalGet(key string) ([]byte, error)
}

// tips：下述是go中常用的将函数转换为接口的方式

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) LocalGet(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

type Group struct {
	name      string
	mainCache *Cache

	// load data(not in mainCache)
	localGetter LocalGetter // load data from local
	peerPicker  PeerPicker  // load data from other server
	// 如何选择远端服务呢，一般采用一致性hash的方法，通过key确定远端服务地址
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func NewGroup(name string, maxBytes int64, getter LocalGetter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	group, ok := groups[name]
	if ok {
		panic("already exists group")
	}
	group = &Group{
		name:        name,
		localGetter: getter,
		mainCache:   NewCache(maxBytes),
	}
	// add group to global groups
	groups[name] = group

	return group
}

func (g *Group) RegisterPicker(picker PeerPicker) {
	if g.peerPicker != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peerPicker = picker
}

func (g *Group) Get(key string) (Byteview, error) {
	if key == "" {
		return Byteview{}, nil
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Printf("key[%s] is hit", key)
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value Byteview, err error) {
	if g.peerPicker != nil {
		if getter, ok := g.peerPicker.PickPeer(key); ok {
			if value, err = g.loadFromPeer(getter, key); err == nil {
				log.Printf("key[%s] get from peer success", key)
				return value, nil
			}
		}
	}

	return g.loadFromLocal(key)
}

func (g *Group) loadFromLocal(key string) (Byteview, error) {
	v, err := g.localGetter.LocalGet(key)
	if err != nil {
		log.Printf("key[%s] get from local failed: err=%v", key, err)
		return Byteview{}, err
	}
	log.Printf("key[%s] get from local success", key)
	value := Byteview{b: helper.SliceCopy(v)}
	g.populate(key, value)
	return value, nil
}

func (g *Group) populate(key string, value Byteview) {
	g.mainCache.set(key, value)
}

func (g *Group) loadFromPeer(getter PeerGetter, key string) (Byteview, error) {
	// tips：保障远端有相同名称的group
	bytes, err := getter.PeerGet(g.name, key)
	if err != nil {
		return Byteview{}, err
	}
	return Byteview{b: bytes}, nil
}
