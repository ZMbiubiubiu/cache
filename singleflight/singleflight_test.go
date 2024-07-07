package singleflight

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSingleflightSerial(t *testing.T) {
	g := SingleFlight{}
	var n int32
	for i := 0; i < 100; i++ {
		g.Do("testKey", func() (interface{}, error) {
			atomic.AddInt32(&n, 1)
			return "", nil
		})
	}
	assert.Equal(t, 100, int(n))
}

func TestSingleflightConcurrent(t *testing.T) {
	g := SingleFlight{}
	var n int32
	for i := 0; i < 30; i++ {
		i := i
		go g.Do("testKey", func() (interface{}, error) {
			fmt.Println("do ", i)
			time.Sleep(10 * time.Second)
			atomic.AddInt32(&n, 1)
			return "", nil
		})
	}
	time.Sleep(11 * time.Second)
	assert.Equal(t, 1, int(n))
}
