package workers

import (
	"encoding/json"
	"time"

	"github.com/syumai/workers/cloudflare/kv"
)

var namespace *kv.Namespace

type Cache struct {
	Content map[string]string `json:"content"`
	Expire  time.Time         `json:"Expire"`
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

func SetCache(key string, content map[string]string) error {
	cache := Cache{
		Content: content,
		Expire:  time.Now().UTC().Add(1 * time.Hour),
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

func GetCache(key string) (Cache, bool, error) {
	value, err := namespace.GetString(key, nil)
	if err != nil {
		return Cache{}, false, err
	}

	if value == "<null>" || value == "" {
		return Cache{}, false, nil
	}

	var cache Cache
	err = json.Unmarshal([]byte(value), &cache)
	if err != nil {
		return Cache{}, false, err
	}

	return cache, true, nil
}

func CacheFunc(key string, handler func() (map[string]string, error)) (map[string]string, error) {
	cache, exists, err := GetCache(key)
	if err != nil {
		return nil, err
	}

	if exists && !cache.IsExpired() {
		return cache.Content, nil
	}

	content, err := handler()
	if err != nil {
		return nil, err
	}

	err = SetCache(key, content)
	if err != nil {
		return nil, err
	}

	return content, nil
}
