package cache

import (
	"strconv"
	"testing"
	"time"
)

func TestLruCache(t *testing.T) {
	counter := 0

	fc := NewFileCache(".")
	if fc == nil {
		t.Fatal("Could not create file cache")
	}

	fetcher := func(key string) (interface{}, *time.Time, error) {
		counter++
		return "fetched " + strconv.Itoa(counter) + " " + string(key), nil, nil
	}

	cc := NewLRUCache(fc, time.Minute, 2, fetcher)

	t.Run("FetchNew", func(t *testing.T) {
		value, expiration, err := cc.Get("hello")
		if err != nil || expiration != nil || value.(string) != "fetched 1 hello" {
			t.Error("Got unexpected results", value, expiration, err)
		}

		// fetch again to ensure same result -- fetcher returns new results each time
		value, expiration, err = cc.Get("hello")
		if err != nil || expiration != nil || value.(string) != "fetched 1 hello" {
			t.Error("Got unexpected results 2nd time", value, expiration, err)
		}
	})

	t.Run("FetchNew", func(t *testing.T) {
		value, expiration, err := cc.Get("hello")
		if err != nil || expiration != nil || value.(string) != "fetched 1 hello" {
			t.Error("Got unexpected results", value, expiration, err)
		}

		value, expiration, err = cc.Get("hello2")
		if err != nil || expiration != nil || value.(string) != "fetched 2 hello2" {
			t.Error("Got unexpected results 2nd time", value, expiration, err)
		}

		value, expiration, err = cc.Get("hello3")
		if err != nil || expiration != nil || value.(string) != "fetched 3 hello3" {
			t.Error("Got unexpected results 3rd time", value, expiration, err)
		}

		fc.Add("hello", nil, nil)
		value, expiration, err = cc.Get("hello")
		if err != nil || expiration != nil || value.(string) != "fetched 4 hello" {
			t.Error("Got unexpected results", value, expiration, err)
		}

		fc.Add("hello", nil, nil)
		fc.Add("hello2", nil, nil)
		fc.Add("hello3", nil, nil)
	})
}
