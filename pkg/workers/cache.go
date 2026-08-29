package workers

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/syumai/workers/cloudflare/kv"
)

var namespace *kv.Namespace

type Cache struct {
	Cached map[string]string `json:"cached"`
	Expire time.Time         `json:"Expire"`
}

func (cache Cache) IsExpired() bool {
	return time.Now().UTC().After(cache.Expire)
}

func init() {
	var err error
	namespace, err = kv.NewNamespace("cache")
	if err != nil {
		panic(err)
	}
}

func SetCache(key string, cached map[string]string, ttl time.Duration) error {
	cache := Cache{
		Cached: cached,
		Expire: time.Now().UTC().Add(ttl),
	}
	value, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	err = namespace.PutString(key, string(value), nil)
	if err != nil {
		return err
	}

	return nil
}

func GetCache(key string) (Cache, error) {
	value, err := namespace.GetString(key, nil)
	if err != nil {
		return Cache{}, err
	}

	if value == "<null>" || value == "" {
		return Cache{}, nil
	}

	var cache Cache
	err = json.Unmarshal([]byte(value), &cache)
	if err != nil {
		return Cache{}, err
	}

	return cache, nil
}

func CacheFunc(key string, handler func() (map[string]string, time.Duration, error)) (map[string]string, error) {
	cache, err := GetCache(key)
	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(cache).IsZero() && !cache.IsExpired() {
		return cache.Cached, nil
	}

	cached, ttl, err := handler()
	if err != nil {
		return nil, err
	}

	err = SetCache(key, cached, ttl)
	if err != nil {
		return nil, err
	}

	return cached, nil
}
