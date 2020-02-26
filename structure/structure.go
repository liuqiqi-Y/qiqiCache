package structure

import (
	"container/list"
	"errors"
)

// Cache 缓存结构体
type Cache struct {
	maxBytes  int64
	nowBytes  int64
	dl        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

// Value 缓存的值需要实现该接口
type Value interface {
	Len() int64
}
type element struct {
	key string
	val Value
}

// New 新建一个Cache
func New(maxbyte int64, callback func(string, Value)) (*Cache, error) {
	if maxbyte <= 0 {
		err := errors.New("maxbyte can not lesser than 0")
		return nil, err
	}
	return &Cache{
		maxBytes:  maxbyte,
		dl:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: callback,
	}, nil
}

// Add 向缓存中添加键值对
func (c *Cache) Add(key string, val Value) error {
	if key == "" || val == nil {
		err := errors.New("key or value is null")
		return err
	}
	if ele, ok := c.cache[key]; ok {
		c.dl.MoveToBack(ele)
		kv, ok := ele.Value.(*element)
		if !ok {
			return errors.New("list value type is not right")
		}
		kv.val = val
		c.nowBytes += (val.Len() - kv.val.Len())
		return nil
	}
	ele := c.dl.PushBack(&element{key, val})
	c.cache[key] = ele
	c.nowBytes += (int64(len(key)) + val.Len())
	for c.maxBytes < c.nowBytes {
		b := c.RemoveOldest()
		if b == false {
			err := errors.New("delete oldest value failed")
			return err
		}
	}
	return nil
}

// Get 返回一个缓存值
func (c *Cache) Get(key string) (Value, error) {
	if key == "" {
		err := errors.New("key is null")
		return nil, err
	}
	if ele, ok := c.cache[key]; ok {
		c.dl.MoveToBack(ele)
		kv, _ := ele.Value.(*element)
		return kv.val, nil
	}
	err := errors.New("has no this key")
	return nil, err
}

// RemoveOldest 移除一个最旧的缓存
func (c *Cache) RemoveOldest() bool {
	ele := c.dl.Front()
	if ele == nil {
		return false
	}
	e := c.dl.Remove(ele)
	kv, _ := e.(*element)
	delete(c.cache, kv.key)
	c.nowBytes -= (int64(len(kv.key)) + kv.val.Len())
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.val)
	}
	return true
}

// Len 返回缓存个数
func (c *Cache) Len() int {
	return c.dl.Len()
}
