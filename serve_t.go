package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/liuqiqi-Y/qiqiCache/concurrency"
	"github.com/liuqiqi-Y/qiqiCache/httppool"
	"github.com/liuqiqi-Y/qiqiCache/serve"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *concurrency.Group {
	getter := concurrency.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		})
	return concurrency.NewGroup("score", getter, 2<<10)
}
func startCacheServer(addr string, addrs []string, g concurrency.Group) {
	peers := httppool.NewHTTPPool(addr)
	peers.Set(addrs...)
	g.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], &serve.Server))
}
