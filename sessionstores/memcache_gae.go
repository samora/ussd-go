package sessionstores

import (
	"encoding/json"
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
)

// MemcacheGAE is memcache session store for Google App Engine
type MemcacheGAE struct {
	ctx context.Context
}

func NewMemcacheGAE(ctx context.Context) *MemcacheGAE {
	return &MemcacheGAE{
		ctx: ctx,
	}
}

func (m MemcacheGAE) Connect() error {
	return nil
}

func (m MemcacheGAE) Close() error {
	return nil
}

func (m MemcacheGAE) SetValue(key, value string) error {
	return memcache.Set(m.ctx, &memcache.Item{
		Key:   key,
		Value: []byte(value),
	})
}

func (m MemcacheGAE) GetValue(key string) (string, error) {
	item, err := memcache.Get(m.ctx, key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

func (m MemcacheGAE) ValueExists(key string) bool {
	_, err := memcache.Get(m.ctx, key)
	if err != nil {
		return false
	}
	return true
}

func (m MemcacheGAE) DeleteValue(key string) error {
	return memcache.Delete(m.ctx, key)
}

func (m MemcacheGAE) HashSetValue(name, key, value string) error {
	hash := map[string]string{}
	_, _ = memcacheGAEHashCodec.Get(m.ctx, name, hash)
	hash[key] = value
	return memcacheGAEHashCodec.Set(m.ctx, &memcache.Item{
		Key:    name,
		Object: hash,
	})
}

func (m MemcacheGAE) HashGetValue(name, key string) (string, error) {
	hash := map[string]string{}
	_, err := memcacheGAEHashCodec.Get(m.ctx, name, hash)
	if err != nil {
		return "", err
	}
	if val, ok := hash[key]; ok {
		return val, nil
	}
	return "", errors.New("Hash value not found")
}

func (m MemcacheGAE) HashValueExists(name, key string) bool {
	hash := map[string]string{}
	_, _ = memcacheGAEHashCodec.Get(m.ctx, name, hash)
	_, ok := hash[key]
	return ok
}

func (m MemcacheGAE) HashDeleteValue(name, key string) error {
	hash := map[string]string{}
	_, _ = memcacheGAEHashCodec.Get(m.ctx, name, hash)
	delete(hash, key)
	return memcacheGAEHashCodec.Set(m.ctx, &memcache.Item{
		Key:    name,
		Object: hash,
	})
}

func (m MemcacheGAE) HashExists(name string) bool {
	return m.ValueExists(name)
}

func (m MemcacheGAE) HashDelete(name string) error {
	return memcache.Delete(m.ctx, name)
}

var memcacheGAEHashCodec = memcache.Codec{
	Marshal: func(v interface{}) ([]byte, error) {
		hash, ok := v.(map[string]string)
		if !ok {
			return []byte{}, errors.New("Invalid type. Expected map[string]string")
		}
		return json.Marshal(&hash)
	},
	Unmarshal: func(b []byte, v interface{}) error {
		hash, ok := v.(map[string]string)
		if !ok {
			return errors.New("Invalid type. Expected map[string]string")
		}
		return json.Unmarshal(b, &hash)
	},
}
