package config

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/stretchr/testify/assert"
)

type fakeLogger struct {
	T *testing.T
}

func (f fakeLogger) Printf(format string, args ...interface{}) {
	f.T.Fatalf(format, args)
}

func makeCacheManager(t *testing.T, logger Logger) (*bigcache.BigCache, *Cache) {
	hasher, err := newHasher()
	if err != nil {
		t.Fatalf("Cannot create hash function for cache: %v", err)
	}

	cc := bigcache.DefaultConfig(time.Minute)
	cc.HardMaxCacheSize = 10
	cc.Hasher = hasher
	cache, err := bigcache.NewBigCache(cc)
	if err != nil {
		t.Fatalf("Cannot create cache: %v", err)
	}
	return cache, &Cache{Cache: cache, Logger: logger}
}

type fakeStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGetCache(t *testing.T) {
	c, cm := makeCacheManager(t, fakeLogger{t})
	expected := fakeStruct{"Alice", 24}
	b, _ := json.Marshal(expected)
	c.Set("alice", b)

	var actual fakeStruct
	cm.Get("alice", &actual)
	assert.Equal(t, expected, actual)
}

func TestSetCache(t *testing.T) {
	c, cm := makeCacheManager(t, fakeLogger{t})
	expected := fakeStruct{"Alice", 24}
	cm.Set("alice", expected)
	b, err := c.Get("alice")
	assert.NoError(t, err)

	var actual fakeStruct
	json.Unmarshal(b, &actual)

	assert.Equal(t, expected, actual)
}

func TestDelCache(t *testing.T) {
	c, cm := makeCacheManager(t, fakeLogger{t})
	c.Set("alice", []byte(`{"name":"Alice","age":24}`))
	cm.Del("alice")

	var actual fakeStruct
	has := cm.Get("alice", &actual)

	assert.False(t, has)
}

func TestFullScenario(t *testing.T) {
	_, cm := makeCacheManager(t, fakeLogger{t})
	expected := fakeStruct{"Alice", 24}
	cm.Set("alice", expected)

	var actual fakeStruct
	has := cm.Get("alice", &actual)
	assert.True(t, has)
	assert.Equal(t, expected, actual)

	cm.Del("alice")

	has = cm.Get("alice", &actual)

	assert.False(t, has)
}
