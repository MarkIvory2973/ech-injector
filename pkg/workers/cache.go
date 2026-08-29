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

func SetCache(key string, cached map[string]string, ttl int) error {
	value, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	err = namespace.PutString(key, string(value), &kv.PutOptions{
		ExpirationTTL: ttl,
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

	var cached map[string]string
	err = json.Unmarshal([]byte(value), &cached)
	if err != nil {
		return nil, err
	}

	return cached, nil
}

func CacheFunc(key string, handler func() (map[string]string, int, error)) (map[string]string, error) {
	cached, err := GetCache(key)
	if err != nil {
		return nil, err
	}

	if cached != nil {
		return cached, nil
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
