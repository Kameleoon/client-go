package storage

import (
	"time"
)

type CacheFactory interface {
	Create(time time.Duration, autoDelete bool) (Cache, error)
}

type CacheFactoryImpl struct {
}

func (cf *CacheFactoryImpl) Create(expirationTime time.Duration, autoDelete bool) (Cache, error) {
	return NewCache(expirationTime, autoDelete)
}
