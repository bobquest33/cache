package cache

import (
	"encoding/json"
	"log"
	"time"
)

type IEncodeDecode interface {
	Encode(v interface{}) ([]byte, error)
	Decode(b []byte) (interface{}, error)
}

type jsonEncoderDecoder struct{}

func (jj *jsonEncoderDecoder) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

func (jj *jsonEncoderDecoder) Decode(b []byte) (interface{}, error) {
	if b == nil {
		return nil, nil
	}
	var result interface{}
	err := json.Unmarshal(b, &result)
	return result, err
}

type marshaler struct {
	cache        ICache
	encodeDecode IEncodeDecode
}

func NewMarshaler(cache ICache, encodeDecode IEncodeDecode) ICache {
	return &marshaler{cache: cache, encodeDecode: encodeDecode}
}

func NewJSONMarshaler(cache ICache) ICache {
	return NewMarshaler(cache, &jsonEncoderDecoder{})
}

func (mm *marshaler) Add(key string, value interface{}, expiration *time.Time) error {
	if val, err := mm.encodeDecode.Encode(value); err != nil {
		return err
	} else {
		return mm.cache.Add(key, val, expiration)
	}
}

func (mm *marshaler) Get(key string) (interface{}, *time.Time, error) {
	if result, expiration, err := mm.cache.Get(key); err != nil || result == nil {
		return nil, nil, err
	} else if bytes, ok := result.([]byte); !ok {
		log.Panic("Unexepcted result type, expected bytes array", result)
		return nil, expiration, err
	} else {
		return bytes, expiration, err
	}
}
