package main

type Cacher interface {
	Read(key string) (string, bool)
	Write(key string, val string)
	Delete(key string)
}

type Cache struct {
	cache   map[string]string
	reads   chan cacheReadRequest
	writes  chan [2]string
	deletes chan string
}

type cacheReadRequest struct {
	key     string
	results chan cacheResponse
}

type cacheResponse struct {
	val string
	ok  bool
}

func NewCache() *Cache {
	c := &Cache{
		cache:   make(map[string]string),
		reads:   make(chan cacheReadRequest),
		writes:  make(chan [2]string),
		deletes: make(chan string),
	}

	go c.listen()
	return c
}

func (c *Cache) Read(key string) (string, bool) {
	responseChan := make(chan cacheResponse)
	c.reads <- cacheReadRequest{key, responseChan}
	response := <-responseChan
	return response.val, response.ok
}

func (c *Cache) Write(key string, val string) {
	c.writes <- [2]string{key, val}
}

func (c *Cache) Delete(key string) {
	c.deletes <- key
}

func (c *Cache) listen() {
	for {
		var request cacheReadRequest
		var tuple [2]string
		var key string
		select {
		case request = <-c.reads:
			val, ok := c.cache[request.key]
			request.results <- cacheResponse{val, ok}
		case tuple = <-c.writes:
			c.cache[tuple[0]] = tuple[1]
		case key = <-c.deletes:
			delete(c.cache, key)
		}
	}
}
