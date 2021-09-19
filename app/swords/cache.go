package swords

import "sync"

type autoRelay struct {
	sync.Mutex
	messages map[string]string
}

var cache autoRelay

func init() {
	cache.messages = make(map[string]string)
}

func (reply *autoRelay) Get(keyword string) string {
	return reply.messages[keyword]
}

func (reply *autoRelay) Set(keyword string, word string) {
	reply.Lock()
	defer reply.Unlock()
	reply.messages[keyword] = word
}
