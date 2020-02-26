package concurrency

import (
	"errors"
	"sync"
)

// Getter 当缓存中没有时从数据源获取
type Getter interface {
	Get(string) ([]byte, error)
}

// GetterFunc 函数签名
type GetterFunc func(string) ([]byte, error)

// Get 接口实现
func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

// Group 一个group包含一个cache
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 新建一个group
func NewGroup(name string, getter Getter, bytes int64) *Group {
	if getter == nil {
		panic("getter is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheSize: bytes},
	}
	groups[name] = g
	return g
}

// GetGroup 返回一个group
func GetGroup(name string) *Group {
	mu.RLock()
	if g, ok := groups[name]; ok {
		return g
	}
	mu.RUnlock()
	return nil
}
func (g *Group) populateCache(key string, value ByteView) error {
	if err := g.mainCache.add(key, value); err != nil {
		return err
	}
	return nil
}
func (g *Group) getlocally(key string) (ByteView, error) {
	if key == "" {
		err := errors.New("key is null")
		return ByteView{}, err
	}
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	err = g.populateCache(key, ByteView{b: bytes})
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

// Get 返回一个缓存值
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is null")
	}
	v, err := g.mainCache.get(key)
	if err != nil {
		return g.getlocally(key)
	}
	return v, nil
}
