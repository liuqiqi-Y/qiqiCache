package structure

import (
	"testing"
)

type String string

func (s String) Len() int64 {
	return int64(len(s))
}
func TestGet(t *testing.T) {
	c, _ := New(2024, nil)
	_ = c.Add("key1", String("123"))
	v, _ := c.Get("key1")
	if string(v.(String)) == "123" {
		t.Logf("success to store: %s", "123")
	} else {
		t.Fatalf("cache hit key1=123 failed")
	}
}
func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := int64(len(k1 + k2 + v1 + v2))
	c, _ := New(cap, nil)
	_ = c.Add(k1, String(v1))
	_ = c.Add(k2, String(v2))
	_ = c.Add(k3, String(v3))
	if _, err := c.Get(k1); err != nil || c.Len() != 2 {
		t.Errorf("wrong: %s", err.Error())
	}
}
