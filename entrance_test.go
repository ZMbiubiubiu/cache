package cache

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = map[string][]byte{
	"1": []byte{'a'},
	"2": []byte{'b'},
	"3": []byte{'c'},
}

var stat = make(map[string]int)

// 模拟从本地读取数据
func MockLocalRead(key string) ([]byte, error) {
	v, ok := data[key]
	if !ok {
		return []byte{}, errors.New("no exists")
	}
	// 每从本地读取一次，进行计数
	stat[key]++
	return v, nil
}

func TestGroup(t *testing.T) {
	group := NewGroup("group1", 0, GetterFunc(MockLocalRead))
	for k := range data {
		value, err := group.Get(k)
		assert.Nil(t, err)
		assert.Equal(t, string(data[k]), value.String())
		assert.Equal(t, 1, stat[k])

		value, err = group.Get(k)
		assert.Nil(t, err)
		assert.Equal(t, string(data[k]), value.String())
		// 再次读取，依然是1次
		assert.Equal(t, 1, stat[k])
	}

	key := "unknown"
	_, err := group.Get(key)
	assert.NotNil(t, err)
}
