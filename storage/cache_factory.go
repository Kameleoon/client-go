package storage

import (
	"time"
)

//  Represents a factory type that is supposed to create an object of some cache type.
type CacheFactory interface {
	// Creates new object of some cache type with the specified expiration period
	// and auto delete option. If it's true then Cache will be automatically clean itself with after timeout
	Create(time time.Duration, autoDelete bool) (Cache, error)
}

type CacheFactoryImpl struct {
}

func (cf *CacheFactoryImpl) Create(expirationTime time.Duration, autoDelete bool) (Cache, error) {
	return NewCache(expirationTime, autoDelete)
}
