package storage

import (
	"container/list"
	"sync"
	"time"
)

// Represents a container type that is supposed to store relevant values by distinct keys
// and to manage relevance of the stored values
type Cache interface {
	// If item stored by specified key exists, this method updates its value and its relevance,
	// otherwise this method creates new relevant item with specified key and value.
	Set(key interface{}, value interface{}, expiration ...time.Duration)
	// Get a value of item by specified key, return an item and bool indicating if value exists
	Get(key string) (interface{}, bool)

	// Removes all stored values.
	Clear()
	// Number of items stored.
	Len() int

	// Returns all keys of stored items (even expired)
	Keys() []interface{}
	// Returns all values of stored items (even expired)
	Values() []interface{}
	// Returns map of actual keys of stored items (only non-expired)
	ActualValues() map[interface{}]interface{}
}

type cacheItem struct {
	key        interface{}
	value      interface{}
	expiration time.Time
}

type ErrCacheExpirationTime struct{}

func (e *ErrCacheExpirationTime) Error() string {
	return "setting an expirationTime parameter to zero or less for a cache doesn't make sense " +
		"because the cache items would be deleted immediately upon insertion"
}

type CacheImpl struct {
	sync.Mutex
	list           *list.List
	values         map[interface{}]*list.Element
	expirationTime time.Duration // expirationTime means in what time the item expired (if won't be overwritten)
	autoDelete     bool          // autoDelete == true runs a goroutine to remove expirated items
}

func NewCache(expirationTime time.Duration, autoDelete bool) (*CacheImpl, error) {
	if expirationTime <= 0 {
		return nil, &ErrCacheExpirationTime{}
	}
	return &CacheImpl{
		list:           list.New(),
		values:         make(map[interface{}]*list.Element),
		expirationTime: expirationTime,
		autoDelete:     autoDelete}, nil
}

// expiration parameter overwrites the default cache expiration time
func (c *CacheImpl) Set(key interface{}, value interface{}, expirationTime ...time.Duration) {
	c.Lock()
	defer c.Unlock()
	expTime := c.expirationTime
	if len(expirationTime) > 0 {
		expTime = expirationTime[0]
	}
	expiration := time.Now().Add(expTime)
	if elem, ok := c.values[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*cacheItem).value = value
		elem.Value.(*cacheItem).expiration = expiration
	} else {
		item := &cacheItem{
			key:        key,
			value:      value,
			expiration: expiration,
		}
		elem := c.list.PushFront(item)
		c.values[key] = elem
		if c.autoDelete && (len(c.values) == 1) {
			// Start the cleanup process if this is the first element
			go c.startCleanup(c.expirationTime)
		}
	}
}

func (c *CacheImpl) Get(key string) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	if elem, ok := c.values[key]; ok {
		item := elem.Value.(*cacheItem)
		if item.expiration.After(time.Now()) {
			return item.value, true
		}
		c.list.Remove(elem)
		delete(c.values, item.key)
	}
	return nil, false
}

// Return all stored keys (even expired)
func (c *CacheImpl) Keys() []interface{} {
	c.Lock()
	defer c.Unlock()
	keys := make([]interface{}, len(c.values))
	i := 0
	for key := range c.values {
		keys[i] = key
		i++
	}
	return keys
}

// Return all stored values (even expired)
func (c *CacheImpl) Values() []interface{} {
	c.Lock()
	defer c.Unlock()
	values := make([]interface{}, len(c.values))
	i := 0
	for _, v := range c.values {
		values[i] = v.Value.(*cacheItem).value
		i++
	}
	return values
}

// Returns a map with actual values
func (c *CacheImpl) ActualValues() map[interface{}]interface{} {
	c.purge()

	c.Lock()
	defer c.Unlock()

	timeNow := time.Now()
	mapValues := make(map[interface{}]interface{}, len(c.values))
	for k, v := range c.values {
		cacheItem := v.Value.(*cacheItem)
		if cacheItem.expiration.After(timeNow) {
			mapValues[k] = cacheItem.value
		} else {
			c.list.Remove(v)
			delete(c.values, k)
		}
	}
	return mapValues
}

func (c *CacheImpl) Len() int {
	return len(c.values)
}

func (c *CacheImpl) Clear() {
	c.Lock()
	defer c.Unlock()

	c.values = make(map[interface{}]*list.Element)
	c.list.Init()
}

// Start ticker for cleainig the data
func (c *CacheImpl) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.purge()
		if c.list.Len() == 0 {
			// Stop the timer if the cache is empty
			return
		}
	}
}

// Remove the expired items from cache
func (c *CacheImpl) purge() {
	c.Lock()
	defer c.Unlock()
	timeNow := time.Now()
	for c.list.Len() > 0 {
		elem := c.list.Back()
		item := elem.Value.(*cacheItem)
		if item.expiration.After(timeNow) {
			break
		}
		c.list.Remove(elem)
		delete(c.values, item.key)
	}
}
