package concurrency

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("12345")
	if v, _ := f.Get("12345"); !reflect.DeepEqual(v, expect) {
		t.Fatalf("getter failed")
	}
}

var db = map[string]string{
	"Tom":   "123",
	"Jac":   "456",
	"Harry": "789",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	var getter Getter = GetterFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key]++
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s is not exist", key)
	})
	g := NewGroup("score", getter, 2<<10)
	for k, v := range db {
		if val, err := g.Get(k); err != nil || val.String() != v {
			t.Fatalf("failed to get value of %s\n", k)
		}
		if _, err := g.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache miss: %s\n", k)
		}
	}
	if view, err := g.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
