package server

import "sync"

type Result struct {
	data map[string]string
	mu   sync.Mutex
}

func (r *Result) Store(key, val string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = val
}

func (r *Result) GetData() map[string]string {
	return r.data
}
