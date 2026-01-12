package main

import (
	"sync"
	"time"
)

type Cache struct {
	feriados   []Feriado
	lastUpdate time.Time
	ttl        time.Duration
	mu         sync.RWMutex
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		ttl: ttl,
	}
}

func (c *Cache) Get() ([]Feriado, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if time.Since(c.lastUpdate) < c.ttl && len(c.feriados) > 0 {
		return c.feriados, true
	}

	return nil, false
}

func (c *Cache) Set(feriados []Feriado) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.feriados = feriados
	c.lastUpdate = time.Now()
}

func (c *Cache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.feriados = nil
	c.lastUpdate = time.Time{}
}
