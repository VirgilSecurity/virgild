package coreapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	thttp "github.com/stretchr/testify/http"
)

func TestWrapperAPIHandlerServeHTTP_SeccessStruct(t *testing.T) {
	seccessObject := map[string]string{
		"test_key": "test_val",
	}
	handler := func(req *http.Request) (interface{}, error) {
		return seccessObject, nil
	}
	l := new(fakeLogger)
	w := &thttp.TestResponseWriter{}
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	wrap := wrapAPIHandler(l)(handler)
	wrap.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.StatusCode)
	assert.Equal(t, `{"test_key":"test_val"}`, w.Output)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestWrapperAPIHandlerServeHTTP_SeccessNil(t *testing.T) {
	handler := func(req *http.Request) (interface{}, error) {
		return nil, nil
	}
	l := new(fakeLogger)
	w := &thttp.TestResponseWriter{}
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	wrap := wrapAPIHandler(l)(handler)
	wrap.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.StatusCode)
	assert.Equal(t, ``, w.Output)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestWrapperAPIHandlerServeHTTP_SeccessByte(t *testing.T) {
	handler := func(req *http.Request) (interface{}, error) {
		return []byte("seccess"), nil
	}
	l := new(fakeLogger)
	w := &thttp.TestResponseWriter{}
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	wrap := wrapAPIHandler(l)(handler)
	wrap.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.StatusCode)
	assert.Equal(t, `seccess`, w.Output)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestWrapperAPIHandlerServeHTTP_APIError(t *testing.T) {
	const code = -1
	const statusCode = 230
	handler := func(req *http.Request) (interface{}, error) {
		return nil, errors.Wrap(APIError{
			Code:       code,
			StatusCode: statusCode,
		}, "")
	}
	l := new(fakeLogger)
	w := &thttp.TestResponseWriter{}
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	wrap := wrapAPIHandler(l)(handler)
	wrap.ServeHTTP(w, r)

	assert.Equal(t, statusCode, w.StatusCode)
	assert.Equal(t, `{"code":-1}`, w.Output)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestWrapperAPIHandlerServeHTTP_InternalErr(t *testing.T) {
	handler := func(req *http.Request) (interface{}, error) {
		return nil, errors.New("Error")
	}
	l := new(fakeLogger)
	l.On("Err").Once()
	w := &thttp.TestResponseWriter{}
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	wrap := wrapAPIHandler(l)(handler)
	wrap.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.StatusCode)
	assert.Equal(t, `{"code":10000}`, w.Output)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestWrapperAPIHandlerServeHTTP_InternalErr_LogIt(t *testing.T) {
	handler := func(req *http.Request) (interface{}, error) {
		return nil, errors.New("Error")
	}
	l := new(fakeLogger)
	l.On("Err").Once()
	w := &thttp.TestResponseWriter{}
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	wrap := wrapAPIHandler(l)(handler)
	wrap.ServeHTTP(w, r)

	l.AssertExpectations(t)
}
