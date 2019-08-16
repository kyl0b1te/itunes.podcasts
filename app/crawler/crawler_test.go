package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>Test Page</title>
</head>
<body>
<a class="target" href="http://x.com/podcasts-test1-first/id1">link #1</a>
<a class="target" href="http://x.com/podcasts-test1-second/id2">link #2</a>
<a class="target" href="http://x.com/podcasts-test2-first/id3">link #3</a>
</body>
</html>
		`))
	})

	mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(404)
		w.Write([]byte("<p>error</p>"))
	})

	return httptest.NewServer(mux)
}

func getMockedEntities() map[string]string {

	return map[string]string{
		"link #1": "http://x.com/podcasts-test1-first/id1",
		"link #2": "http://x.com/podcasts-test1-second/id2",
		"link #3": "http://x.com/podcasts-test2-first/id3",
	}
}

func TestGetEntityURLs(t *testing.T) {

	ts := newTestServer()
	defer ts.Close()

	entities, _ := GetEntities(&RequestOptions{
		LookupURL: ts.URL,
		Pattern:   ".target",
	})
	mocked := getMockedEntities()

	assert.Equal(t, len(mocked), len(entities))
	for name, url := range mocked {
		assert.Contains(t, mocked, name)
		assert.Equal(t, url, entities[name])
	}

	_, err := GetEntities(&RequestOptions{
		LookupURL: ts.URL + "/404",
		Pattern:   ".target",
	})
	assert.Equal(t, "Not Found", err.Error())
}

func TestGetEntityIDFromURL(t *testing.T) {

	id, err := GetEntityIDFromURL("http://x.x/a/id1")
	assert.Equal(t, 1, id)
	assert.Nil(t, err)

	id, _ = GetEntityIDFromURL("http://x.x/a/b/c/d/e/id56635645")
	assert.Equal(t, 56635645, id)
	assert.Nil(t, err)

	_, err = GetEntityIDFromURL("http://x.x/a/idd")
	assert.Equal(t, "strconv.Atoi: parsing \"d\": invalid syntax", err.Error())
}
