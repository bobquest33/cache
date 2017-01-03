// package cache implements a very simple map like interface for caches
package cache

import "time"

// ICache defines a simple cache interface.
type ICache interface {
	Add(key string, value interface{}, expiration *time.Time) error
	Get(key string) (interface{}, *time.Time, error)
}
