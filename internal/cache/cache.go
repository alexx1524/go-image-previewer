package cache

import (
	lrucache "github.com/alexx1524/go-home-work/hw04_lru_cache"
)

type Cache interface {
	SetRemoveItemsChan(chan<- lrucache.RemovedItem)
	Set(key lrucache.Key, value interface{}) bool
	Get(key lrucache.Key) (interface{}, bool)
	Clear()
}
