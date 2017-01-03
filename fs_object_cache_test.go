package cache

import (
	"testing"
)

func TestFileCache(t *testing.T) {
	cases := []string{"", "."}
	for _, tc := range cases {
		t.Run("Dir = "+tc, func(t *testing.T) {
			fc := NewFileCache(tc)
			if fc == nil {
				t.Fatal("Could not create file cache")
			}

			t.Run("ReadNotExists", func(t *testing.T) {
				value, expiration, err := fc.Get("not exists")
				if err != nil || value != nil || expiration != nil {
					t.Error("Got unexpected results", value, expiration, err)
				}
			})
			t.Run("WriteSimple", func(t *testing.T) {
				value := "hello world"
				if err := fc.Add("simple", value, nil); err != nil {
					t.Error("Could not write simple value", err)
				}
			})
			t.Run("ReadSimple", func(t *testing.T) {
				value, expiration, err := fc.Get("simple")
				if err != nil || expiration != nil || value != "hello world" {
					t.Error("Could not read simple value", err, expiration, value)
				}
			})
			t.Run("DeleteSimple", func(t *testing.T) {
				if err := fc.Add("simple", nil, nil); err != nil {
					t.Error("Failed to delete simple", err)
				}
			})
		})
	}
}
