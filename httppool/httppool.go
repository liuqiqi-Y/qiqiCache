package httppool

import (
	"log"
	"sync"

	"github.com/liuqiqi-Y/qiqiCache/consistenthash"
	"github.com/liuqiqi-Y/qiqiCache/serve"
)

const (
	defaultReplicas = 50
)

type PeerPicker interface {
	PickPeer(string) (PeerGetter, bool)
}
type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
}

func NewHTTPPool(s string) *HTTPPool {
	return &HTTPPool{
		self:     s,
		basePath: serve.DefaultBasePath,
	}
}
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.peers == nil {
		p.peers = consistenthash.New(defaultReplicas, nil)
	}
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		log.Printf("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)
