package lru

import (
	"container/list"
	"errors"
)

var ErrTooLargeEntry = errors.New("key/value pair is too large")

type Cache struct {
	maxBytes int64 // 最多存储的字节数，若=0表示不设置上限
	useBytes int64 // 使用的字节数
	ll       *list.List
	cache    map[string]*list.Element

	OnEvicted func(key string, value Value) // 移除键值对时的回调函数
}

type Value interface {
	Len() int // 返回该Value占用的字节数
}

type entry struct {
	key   string
	value Value
}

func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	cache := &Cache{
		maxBytes:  maxBytes,
		useBytes:  0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}

	return cache
}

func (c *Cache) Len() int {
	return len(c.cache)
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	ele, ok := c.cache[key]
	if !ok {
		return nil, ok
	}
	kv := ele.Value.(*entry)
	// 将ele插入到链表头部
	c.ll.MoveToFront(ele)
	return kv.value, ok
}

func (c *Cache) Set(key string, value Value) error {
	// 不能设置过大的键值对
	if len(key)+value.Len() > int(c.maxBytes) {
		return ErrTooLargeEntry
	}
	ele, ok := c.cache[key]
	// 本次设置键值对引起的占用字节数变化
	var changeLen = len(key) + value.Len()
	if ok {
		oldValue := ele.Value.(*entry).value
		changeLen = value.Len() - oldValue.Len()
	}

	// 若超过内存占用上限，先采用lru策略淘汰键值对
	for c.maxBytes != 0 && int64(changeLen)+c.useBytes > c.maxBytes {
		ele := c.ll.Back()
		if ele == nil {
			break
		}
		c.Delete(ele.Value.(*entry).key)
	}
	// 设置键值对
	c.useBytes += int64(changeLen)
	var kv *entry
	if !ok {
		kv = &entry{
			key:   key,
			value: value,
		}
		ele = c.ll.PushFront(kv)
	} else {
		kv = ele.Value.(*entry)
		kv.value = value
		c.ll.MoveToFront(ele)
	}
	c.cache[key] = &list.Element{Value: kv}
	return nil
}

func (c *Cache) Delete(key string) {
	ele, ok := c.cache[key]
	if !ok {
		return
	}
	kv := ele.Value.(*entry)

	delete(c.cache, key)
	c.ll.Remove(ele)
	c.useBytes -= int64(len(key) + kv.value.Len())
	if c.OnEvicted != nil {
		c.OnEvicted(key, kv.value)
	}
}
