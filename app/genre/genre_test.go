package genre

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
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

func getMockedGenres() []*Genre {

	return []*Genre{
		NewGenre(1, "link #1", "http://x.com/podcasts-test1-first/id1"),
		NewGenre(2, "link #2", "http://x.com/podcasts-test1-second/id2"),
		NewGenre(3, "link #3", "http://x.com/podcasts-test2-first/id3"),
	}
}

func TestGetGenresFromWeb(t *testing.T) {

	ts := newTestServer()
	defer ts.Close()

	genres, _ := GetGenres(&crawler.RequestOptions{
		LookupURL: ts.URL,
		Pattern:   ".target",
	})
	mocked := getMockedGenres()

	assert.Equal(t, len(mocked), len(genres))
	for _, genre := range mocked {
		assert.Contains(t, genres, genre)
	}

	_, err := GetGenres(&crawler.RequestOptions{
		LookupURL: ts.URL + "/404",
		Pattern:   ".target",
	})
	assert.Equal(t, "Not Found", errors.Cause(err).Error())
}

func TestSaveGenres(t *testing.T) {

	genres := getMockedGenres()

	err := SaveGenres("/tmp/inner/path/genres.json", genres)
	assert.Equal(t, "open /tmp/inner/path/genres.json: no such file or directory", err.Error())

	err = SaveGenres("/tmp/genres.json", genres)
	assert.Nil(t, err)
	assert.FileExists(t, "/tmp/genres.json")

	os.Remove("/tmp/genres.json")
}
