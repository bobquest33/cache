package cache

import (
	"sync"
	"time"
)

type entry struct {
	err           error
	value         interface{}
	expiration    *time.Time
	ttlExpiration *time.Time
}

type LruFetcher func(key string) (interface{}, *time.Time, error)
type lru struct {
	cache ICache

	sync.Mutex
	values map[string]*entry

	ttl      time.Duration
	maxCount int64
	fetcher  LruFetcher
}

func NewLRUCache(cache ICache, ttl time.Duration, maxCount int64, fetcher LruFetcher) ICache {
	return &lru{
		cache:    cache,
		values:   map[string]*entry{},
		ttl:      ttl,
		maxCount: maxCount,
		fetcher:  fetcher,
	}
}

func (cc *lru) Add(key string, value interface{}, expiration *time.Time) error {
	defer cc.Unlock()
	cc.Lock()
	return cc.lockedAdd(key, value, expiration)
}

func (cc *lru) Get(key string) (interface{}, *time.Time, error) {
	defer cc.Unlock()
	defer cc.lockedEnsureMaxCount()
	cc.Lock()
	return cc.lockedGet(key)
}

func (cc *lru) lockedAdd(key string, value interface{}, expiration *time.Time) error {
	delete(cc.values, key)
	if err := cc.cache.Add(key, value, expiration); err != nil {
		return err
	}
	return nil
}

func (cc *lru) lockedGet(key string) (interface{}, *time.Time, error) {
	now := time.Now()
	if ee, ok := cc.values[key]; ok {
		if ee.ttlExpiration.After(now) {
			return ee.value, ee.expiration, ee.err
		}
	}

	value, expiration, err := cc.cache.Get(key)
	if err == nil && value == nil && expiration == nil && cc.fetcher != nil {
		value, expiration, err = cc.fetcher(key)
		if err == nil {
			cc.cache.Add(key, value, expiration)
		}
	}

	ttlExpiration := getExpiration(expiration, cc.ttl, err)
	cc.values[key] = &entry{
		value:         value,
		err:           err,
		expiration:    expiration,
		ttlExpiration: ttlExpiration,
	}
	return value, expiration, err
}

func (cc *lru) lockedEnsureMaxCount() {
	if cc.maxCount > 0 && int64(len(cc.values)) > cc.maxCount {
		var key string
		var value *entry
		for k, v := range cc.values {
			if value == nil || v.ttlExpiration.Before(*value.ttlExpiration) {
				key = k
				value = v
			}
		}
		delete(cc.values, key)
	}
}

func getExpiration(expiration *time.Time, ttl time.Duration, err error) *time.Time {
	if err != nil {
		ttl = time.Second * 10
	}
	candidate := time.Now().Add(ttl)
	if expiration == nil || expiration.After(candidate) {
		return &candidate
	}
	return expiration
}
