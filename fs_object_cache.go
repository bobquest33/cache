package cache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type fileCache struct {
	dir string
}

func NewFileCache(directory string) ICache {
	if directory == "" {
		if dir, err := ioutil.TempDir("", "fsoc"); err != nil {
			return nil
		} else {
			directory = dir
		}
	}
	return &fileCache{dir: directory}
}

func (fc *fileCache) Add(key string, value interface{}, expiration *time.Time) error {
	if value == nil {
		return os.Remove(fc.dir + "/" + key)
	}

	encoded, err := json.Marshal(&struct {
		Value      interface{}
		Expiration *time.Time
	}{value, expiration})

	if err != nil {
		return err
	}
	return ioutil.WriteFile(fc.dir+"/"+key, encoded, 0666)
}

func (fc *fileCache) Get(key string) (value interface{}, expiration *time.Time, err error) {
	var results []byte

	if results, err = ioutil.ReadFile(fc.dir + "/" + key); err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	var decoded struct {
		Value      interface{}
		Expiration *time.Time
	}
	if err = json.Unmarshal(results, &decoded); err == nil {
		value = decoded.Value
		expiration = decoded.Expiration
	}

	if expiration != nil && expiration.Before(time.Now()) {
		expiration = nil
		value = nil
		err = os.Remove(fc.dir + "/" + key)
	}
	return
}
