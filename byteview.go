// 本项目使用的真实缓存值

package cache

import "cache/helper"

type Byteview struct {
	b []byte
}

func (bv Byteview) Len() int {
	return len(bv.b)
}

func (bv Byteview) String() string {
	return string(bv.b)
}

func (bv Byteview) ByteSlice() []byte {
	return helper.SliceCopy(bv.b)
}
