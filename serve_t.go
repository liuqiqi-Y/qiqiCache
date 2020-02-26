package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/liuqiqi-Y/qiqiCache/concurrency"
	"github.com/liuqiqi-Y/qiqiCache/serve"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	_ = concurrency.NewGroup("scores", concurrency.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), 2<<10)
	addr := "localhost:8000"
	peers := serve.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
