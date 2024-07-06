package consistenthash

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsistentHash(t *testing.T) {
	ch := NewConsistentHash(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	ch.AddNodes("2", "4", "6")
	// Given the above hash function and replicas, this will give virtual nodes like below:
	// real node | virtual nodes
	// 		2		2 12 22
	// 		4		4 14 24
	// 		6		6 16 26

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		assert.Equal(t, v, ch.GetNode(k))
	}

	ch.AddNodes("8")

	// 27 should now map to 8.
	// real node | virtual nodes
	// 		2		2 12 22
	// 		4		4 14 24
	// 		6		6 16 26
	//		8		8 18 28
	testCases["27"] = "8"

	for k, v := range testCases {
		assert.Equal(t, v, ch.GetNode(k))
	}
}
