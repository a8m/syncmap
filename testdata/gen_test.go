package main

import (
	"net/http"
	"testing"
)

func TestIntMap(t *testing.T) {
	var m IntMap
	m.Store(1, 2)
	_, ok := m.Load(1)
	if !ok {
		t.Fatal("value should be existed")
	}
	m.Delete(1)
	_, ok = m.Load(1)
	if ok {
		t.Fatal("value should not be existed")
	}
	r, loaded := m.LoadOrStore(1, 2)
	if loaded {
		t.Fatal("value should not be loaded")
	}
	lr, loaded := m.LoadOrStore(1, r)
	if !loaded {
		t.Fatal("value should not be loaded")
	}
	if lr != r {
		t.Fatal("loaded value should be the same")
	}
	s, _ := m.LoadOrStore(2, 3)
	kv := map[int]int{1: r, 2: s}
	m.Range(func(key, value int) bool {
		v, ok := kv[key]
		if !ok {
			t.Fatal("keys do not match")
		}
		if value != v {
			t.Fatal("values do not match")
		}
		delete(kv, key)
		return true
	})
}

func TestRequests(t *testing.T) {
	var m Requests
	m.Store("r", &http.Request{})
	_, ok := m.Load("r")
	if !ok {
		t.Fatal("value should be existed")
	}
	m.Delete("r")
	_, ok = m.Load("r")
	if ok {
		t.Fatal("value should not be existed")
	}
	r, loaded := m.LoadOrStore("r", &http.Request{})
	if loaded {
		t.Fatal("value should not be loaded")
	}
	lr, loaded := m.LoadOrStore("r", r)
	if !loaded {
		t.Fatal("value should not be loaded")
	}
	if lr != r {
		t.Fatal("loaded value should be the same")
	}
	s, _ := m.LoadOrStore("s", &http.Request{})
	kv := map[string]*http.Request{"r": r, "s": s}
	m.Range(func(key string, value *http.Request) bool {
		v, ok := kv[key]
		if !ok {
			t.Fatal("keys do not match")
		}
		if value != v {
			t.Fatal("values do not match")
		}
		delete(kv, key)
		return true
	})
}
