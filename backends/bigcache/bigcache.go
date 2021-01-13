package bigcache

import (
	"time"

	"github.com/allegro/bigcache"
)

// Bigcache resolver.
type Bigcache struct {
	client *bigcache.BigCache
}

// New returns a new Bigcache instance.
func New() *Bigcache {
	client, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		panic(err)
	}
	return &Bigcache{
		client: client,
	}
}

// Get reads entry for the key.
func (b *Bigcache) Get(key string) ([]byte, error) {
	return b.client.Get(key)
}

// Contains checks existence of the key.
func (b *Bigcache) Contains(key string) (bool, error) {
	panic("not implemented") // TODO: Implement
}

// Set sets the value for the key.
func (b *Bigcache) Set(key string, value []byte) error {
	return b.client.Set(key, value)
}

// Del deletes the key and its associated value.
func (b *Bigcache) Del(key string) error {
	return b.client.Delete(key)
}

// Purge resets the cached keys.
func (b *Bigcache) Purge() error {
	return b.client.Reset()
}

// Len returns the total number of currently stored items.
func (b *Bigcache) Len() (int, error) {
	return b.client.Len(), nil
}
