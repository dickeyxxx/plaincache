package main

import "fmt"

type Cacher interface {
	Read(key string) (string, bool)
	Write(key string, val string)
	Delete(key string)
}

type Cache struct {
	cache   map[string]string
	writes  chan [2]string
	deletes chan string
}

func NewCache() *Cache {
	c := &Cache{
		cache:   make(map[string]string),
		writes:  make(chan [2]string),
		deletes: make(chan string),
	}

	go c.listen()
	return c
}

func (c *Cache) Read(key string) (string, bool) {
	val, ok := c.cache[key]
	return val, ok
}

func (c *Cache) Write(key string, val string) {
	c.writes <- [2]string{key, val}
}

func (c *Cache) Delete(key string) {
	fmt.Println("deleting", key)
	c.deletes <- key
	fmt.Println("deleted", key)
	delete(c.cache, key)
}

func (c *Cache) listen() {
	for {
		var tuple [2]string
		var key string
		select {
		case tuple = <-c.writes:
			c.cache[tuple[0]] = tuple[1]
		case key = <-c.deletes:
			fmt.Println("deleting in goroutine", key)
			delete(c.cache, key)
		}
	}
}
