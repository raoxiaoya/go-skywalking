package util

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
)

type GoroutineLocalStorage struct {
	sync.Map
}

func (gls *GoroutineLocalStorage) Get(key uint64) interface{} {
	value, ok := gls.Load(key)
	if !ok {
		return nil
	}
	return value
}

func (gls *GoroutineLocalStorage) Set(key uint64, v interface{}) {
	gls.Store(key, v)
}

func (gls *GoroutineLocalStorage) Del(key uint64) {
	gls.Delete(key)
}

func (gls *GoroutineLocalStorage) GetGoroutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
