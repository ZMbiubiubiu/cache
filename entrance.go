package cache

import (
	"log"
	"sync"

	"cache/helper"
)

// Getter 用户自己实现各自加载数据功能，用于查找缓存失败时从本地加载
type Getter interface {
	Get(key string) ([]byte, error)
}

// tips：下述是go中常用的将函数转换为接口的方式

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

type Group struct {
	name      string
	getter    Getter
	mainCache *Cache
}

func NewGroup(name string, maxBytes int64, getter Getter) *Group {
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
		name:      name,
		getter:    getter,
		mainCache: NewCache(maxBytes),
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
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

func (g *Group) load(key string) (Byteview, error) {
	return g.loadFromLocal(key)
}

func (g *Group) loadFromLocal(key string) (Byteview, error) {
	v, err := g.getter.Get(key)
	if err != nil {
		return Byteview{}, err
	}
	value := Byteview{b: helper.SliceCopy(v)}
	g.populate(key, value)
	return value, nil
}

func (g *Group) populate(key string, value Byteview) {
	g.mainCache.set(key, value)
}
