package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32
type Map struct {
	hash     Hash
	replicas int
	nodes    []int
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}
func (m *Map) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < m.replicas; i++ {
			h := int(m.hash([]byte(strconv.Itoa(i) + node)))
			m.nodes = append(m.nodes, h)
			m.hashMap[h] = node
		}
	}
	sort.Ints(m.nodes)
}
func (m *Map) Get(key string) string {
	if len(m.nodes) == 0 {
		return ""
	}
	h := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.nodes), func(i int) bool {
		return m.nodes[i] >= h
	})
	return m.hashMap[m.nodes[idx%len(m.nodes)]]
}
