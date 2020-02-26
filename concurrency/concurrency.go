package concurrency

import (
	"fmt"
	"sync"

	"github.com/liuqiqi-Y/qiqiCache/structure"
)

// ByteView 缓存的值
type ByteView struct {
	b []byte
}

// Len 值的长度
func (v ByteView) Len() int64 {
	return int64(len(v.b))
}

// String 值的字符串输出
func (v ByteView) String() string {
	return string(v.b)
}

// ByteSlice 返回值的切片输出
func (v ByteView) ByteSlice() []byte {
	c := make([]byte, len(v.b))
	copy(c, v.b)
	return c
}

type cache struct {
	mu        sync.Mutex
	cc        *structure.Cache
	cacheSize int64
}

func (c *cache) add(key string, value ByteView) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cc == nil {
		ci, err := structure.New(c.cacheSize, nil)
		if err != nil {
			return err
		}
		c.cc = ci
	}
	err := c.cc.Add(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) get(key string) (ByteView, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cc == nil {
		return ByteView{}, fmt.Errorf("no cache struct")
	}
	v, err := c.cc.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	return v.(ByteView), nil
}
