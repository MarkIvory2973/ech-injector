package workers

import (
	"encoding/json"

	"github.com/syumai/workers/cloudflare/kv"
)

var namespace *kv.Namespace

func init() {
	var err error
	namespace, err = kv.NewNamespace("cache")
	if err != nil {
		panic(err)
	}
}

func SetCache(key string, cache map[string]string, expirationTTL int) error {
	value, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	err = namespace.PutString(key, string(value), &kv.PutOptions{
		ExpirationTTL: expirationTTL,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetCache(key string) (map[string]string, error) {
	value, err := namespace.GetString(key, nil)
	if err != nil {
		return nil, err
	}

	if value == "<null>" || value == "" {
		return nil, nil
	}

	var cache map[string]string
	err = json.Unmarshal([]byte(value), &cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func CacheFunc(key string, handler func() (map[string]string, int, error)) (map[string]string, error) {
	cache, err := GetCache(key)
	if err != nil {
		return nil, err
	}

	if len(cache) != 0 {
		return cache, nil
	}

	cache, expirationTTL, err := handler()
	if err != nil {
		return nil, err
	}

	err = SetCache(key, cache, expirationTTL)
	if err != nil {
		return nil, err
	}

	return cache, nil
}
