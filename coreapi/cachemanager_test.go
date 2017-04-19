package coreapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeLogger struct {
	mock.Mock
}

func (f *fakeLogger) Info(format string, args ...interface{}) {
	f.Called()
}
func (f *fakeLogger) Warn(format string, args ...interface{}) {
	f.Called()
}
func (f *fakeLogger) Err(format string, args ...interface{}) {
	f.Called()
}

type fakeCache struct {
	mock.Mock
}

func (f *fakeCache) Get(key string, val interface{}) (bool, error) {
	args := f.Called(key, val)
	return args.Bool(0), args.Error(1)
}

func (f *fakeCache) Set(key string, val interface{}) error {
	args := f.Called(key, val)
	return args.Error(0)
}

func (f *fakeCache) Del(key string) error {
	args := f.Called(key)
	return args.Error(0)
}

func TestCacheManagerGet_CorrectWork(t *testing.T) {

	table := map[string]bool{
		"exist_key":     true,
		"not_exist_key": false,
	}

	for key, val := range table {
		l := new(fakeLogger)
		c := new(fakeCache)
		c.On("Get", key, mock.Anything).Return(val, nil)
		cm := cacheManager{
			logger: l,
			cache:  c,
		}
		has := cm.Get(key, nil)

		assert.Equal(t, val, has)
	}
}

func TestCacheManagerGet_InternalErr_LogErr(t *testing.T) {
	l := new(fakeLogger)
	l.On("Err").Once()
	c := new(fakeCache)
	c.On("Get", "test_key", mock.Anything).Return(false, fmt.Errorf("ERROR"))
	cm := cacheManager{
		logger: l,
		cache:  c,
	}

	cm.Get("test_key", nil)

	l.AssertCalled(t, "Err")
}

func TestCacheManagerGet_InternalErr_ReturnFalse(t *testing.T) {
	l := new(fakeLogger)
	l.On("Err")
	c := new(fakeCache)
	c.On("Get", "test_key", mock.Anything).Return(true, fmt.Errorf("ERROR"))
	cm := cacheManager{
		logger: l,
		cache:  c,
	}

	has := cm.Get("test_key", nil)

	assert.False(t, has)
}

func TestCacheManagerSet_CallInternalCacheSet(t *testing.T) {
	val := "same value"
	l := new(fakeLogger)
	c := new(fakeCache)
	c.On("Set", "key", val).Return(nil).Once()
	cm := cacheManager{
		logger: l,
		cache:  c,
	}
	cm.Set("key", val)

	c.AssertExpectations(t)
}

func TestCacheManagerSet_InternalErr_LogErr(t *testing.T) {
	l := new(fakeLogger)
	l.On("Err").Once()
	c := new(fakeCache)
	c.On("Set", "key", mock.Anything).Return(fmt.Errorf("ERROR"))
	cm := cacheManager{
		logger: l,
		cache:  c,
	}

	cm.Set("key", nil)

	l.AssertCalled(t, "Err")
}

func TestCacheManagerDel_CallInternalCacheDel(t *testing.T) {
	l := new(fakeLogger)
	c := new(fakeCache)
	c.On("Del", "key").Return(nil).Once()
	cm := cacheManager{
		logger: l,
		cache:  c,
	}
	cm.Del("key")

	c.AssertExpectations(t)
}

func TestCacheManagerDel_InternalErr_LogErr(t *testing.T) {
	l := new(fakeLogger)
	l.On("Err").Once()
	c := new(fakeCache)
	c.On("Del", "key").Return(fmt.Errorf("ERROR"))
	cm := cacheManager{
		logger: l,
		cache:  c,
	}

	cm.Del("key")

	l.AssertCalled(t, "Err")
}

//
// func TestGetCache(t *testing.T) {
// 	c, cm := makeCacheManager(t, fakeLogger{t})
// 	expected := fakeStruct{"Alice", 24}
// 	b, _ := json.Marshal(expected)
// 	c.Set("alice", b)
//
// 	var actual fakeStruct
// 	cm.Get("alice", &actual)
// 	assert.Equal(t, expected, actual)
// }
//
// func TestSetCache(t *testing.T) {
// 	c, cm := makeCacheManager(t, fakeLogger{t})
// 	expected := fakeStruct{"Alice", 24}
// 	cm.Set("alice", expected)
// 	b, err := c.Get("alice")
// 	assert.NoError(t, err)
//
// 	var actual fakeStruct
// 	json.Unmarshal(b, &actual)
//
// 	assert.Equal(t, expected, actual)
// }
//
// func TestDelCache(t *testing.T) {
// 	c, cm := makeCacheManager(t, fakeLogger{t})
// 	c.Set("alice", []byte(`{"name":"Alice","age":24}`))
// 	cm.Del("alice")
//
// 	var actual fakeStruct
// 	has := cm.Get("alice", &actual)
//
// 	assert.False(t, has)
// }
//
// func TestFullScenario(t *testing.T) {
// 	_, cm := makeCacheManager(t, fakeLogger{t})
// 	expected := fakeStruct{"Alice", 24}
// 	cm.Set("alice", expected)
//
// 	var actual fakeStruct
// 	has := cm.Get("alice", &actual)
// 	assert.True(t, has)
// 	assert.Equal(t, expected, actual)
//
// 	cm.Del("alice")
//
// 	has = cm.Get("alice", &actual)
//
// 	assert.False(t, has)
// }
