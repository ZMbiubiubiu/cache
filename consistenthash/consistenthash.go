// 一致性hash模块

package consistenthash

import (
	"fmt"
	"hash/crc32"
	"sort"
)

type Hash func(data []byte) uint32

type ConsistentHash struct {
	hash            Hash  // 采用依赖注入的方式，允许用户替换，比如方便测试的场景
	replicas        int   // 每个节点的副本数（虚拟节点数）
	virtualNodeHash []int // 存储虚拟节点的hash值切片
	virtual2Real    map[int]string
}

func NewConsistentHash(replicas int, hash Hash) *ConsistentHash {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}
	return &ConsistentHash{
		hash:            hash,
		replicas:        replicas,
		virtualNodeHash: nil,
		virtual2Real:    make(map[int]string),
	}
}

func (c *ConsistentHash) AddNodes(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < c.replicas; i++ {
			// 0node,1node,2node....
			virtualNode := fmt.Sprintf("%d%s", i, node)
			h := int(c.hash([]byte(virtualNode)))
			c.virtualNodeHash = append(c.virtualNodeHash, h)
			c.virtual2Real[h] = node
		}
	}
	sort.Ints(c.virtualNodeHash)
}

func (c *ConsistentHash) GetNode(key string) (node string) {
	if key == "" {
		return ""
	}

	// 找到key在hash环上的最近一个虚拟节点的下标
	h := int(c.hash([]byte(key)))
	idx := sort.Search(len(c.virtualNodeHash), func(i int) bool {
		return c.virtualNodeHash[i] >= h
	})
	idx = idx % len(c.virtualNodeHash)
	vNodeHash := c.virtualNodeHash[idx]

	return c.virtual2Real[vNodeHash]
}
