package cache

import (
	"testing"
)

func TestS3Cache(t *testing.T) {
	cc := NewS3ObjectCache("s3://kr-hacks/cache", "us-west-2")
	if cc == nil {
		t.Fatal("Could not create file cache")
	}

	t.Run("ReadNotExists", func(t *testing.T) {
		value, expiration, err := cc.Get("not exists")
		if err != nil || value != nil || expiration != nil {
			t.Error("Got unexpected results", value, expiration, err)
		}
	})
	t.Run("WriteSimple", func(t *testing.T) {
		value := []byte("hello world")
		if err := cc.Add("simple", value, nil); err != nil {
			t.Error("Could not write simple value", err)
		}
	})
	t.Run("ReadSimple", func(t *testing.T) {
		value, expiration, err := cc.Get("simple")
		if err != nil || expiration != nil || string(value.([]byte)) != "hello world" {
			t.Error("Could not read simple value", err, expiration, value)
		}
	})
	t.Run("DeleteSimple", func(t *testing.T) {
		if err := cc.Add("simple", nil, nil); err != nil {
			t.Error("Failed to delete simple", err)
		}
	})
}
