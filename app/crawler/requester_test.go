package crawler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func newRequesterTestServer() *httptest.Server {
	mux := http.NewServeMux()

	for url, data := range getTestData() {
		func(url string, data []byte) {
			mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/json")
				w.WriteHeader(200)
				w.Write(data)
			})
		}(url, data)
	}

	mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(404)
		w.Write([]byte("<p>error</p>"))
	})

	return httptest.NewServer(mux)
}

func getTestData() map[string][]byte {

	return map[string][]byte{
		"/test/1": []byte(`{"test": 1}`),
		"/test/2": []byte(`{"test": 2}`),
	}
}

func TestRequestEntities(t *testing.T) {

	ts := newRequesterTestServer()
	defer ts.Close()

	tests := getTestData()

	opt := &RequestOptions{LookupURL: []string{}}
	for url, _ := range tests {
		opt.LookupURL = append(opt.LookupURL, ts.URL+url)
	}

	decoder := func(url string, body []byte) (interface{}, error) {
		tBody, ok := tests[strings.ReplaceAll(url, ts.URL, "")]
		assert.True(t, ok)
		assert.Equal(t, body, tBody)

		return body, nil
	}
	results := RequestEntities(opt, decoder)

	for entity := range results {
		tBody, ok := tests[strings.ReplaceAll(entity.URL, ts.URL, "")]
		assert.True(t, ok)
		assert.Nil(t, entity.Error)
		assert.Equal(t, entity.Entity, tBody)
	}

	opt = &RequestOptions{LookupURL: []string{ts.URL + "/404"}}
	results = RequestEntities(opt, func(url string, body []byte) (interface{}, error) {
		return nil, nil
	})

	for en := range results {
		assert.NotNil(t, en.Error)
		assert.Equal(
			t,
			fmt.Sprintf("Unreachable URL: %s", en.URL),
			errors.Cause(en.Error).Error(),
		)
	}
}

func TestRequestEntitiesWithLimiter(t *testing.T) {

	ts := newRequesterTestServer()
	defer ts.Close()

	tests := getTestData()

	opt := &LimitedRequestOptions{LookupURL: []string{}, Duration: time.Second}
	for url, _ := range tests {
		opt.LookupURL = append(opt.LookupURL, ts.URL+url)
	}

	decoder := func(url string, body []byte) (interface{}, error) {
		tBody, ok := tests[strings.ReplaceAll(url, ts.URL, "")]
		assert.True(t, ok)
		assert.Equal(t, body, tBody)

		return body, nil
	}
	results := RequestEntitiesWithLimiter(opt, decoder)

	for entity := range results {
		tBody, ok := tests[strings.ReplaceAll(entity.URL, ts.URL, "")]
		assert.True(t, ok)
		assert.Nil(t, entity.Error)
		assert.Equal(t, entity.Entity, tBody)
	}

	opt = &LimitedRequestOptions{
		LookupURL: []string{ts.URL + "/404"},
		Duration:  time.Second,
	}
	results = RequestEntitiesWithLimiter(opt, func(url string, body []byte) (interface{}, error) {
		return nil, nil
	})

	for en := range results {
		assert.NotNil(t, en.Error)
		assert.Equal(
			t,
			fmt.Sprintf("Unreachable URL: %s", en.URL),
			errors.Cause(en.Error).Error(),
		)
	}
}
