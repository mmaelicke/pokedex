package internal

import (
	"fmt"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	const interval = time.Second * 5
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "test1",
			val: []byte{},
		},
		{
			key: "test2",
			val: []byte("Long chunk of data"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key %s", c.key)
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value: %v but got %v", c.val, val)
				return
			}
		})
	}
}

func TestEmptyingCache(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond

	cache := NewCache(baseTime)
	cache.Add("foo", []byte("bar"))

	// here the key should be in the cache
	_, ok := cache.Get("foo")
	if !ok {
		t.Errorf("expected to find key foo")
		return
	}

	// wait until it is gone
	time.Sleep(waitTime)

	_, ok = cache.Get("foo")
	if ok {
		t.Errorf("The key foo was not tidy up.")
	}
}
