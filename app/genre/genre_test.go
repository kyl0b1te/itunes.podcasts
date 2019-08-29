package genre

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"

	"github.com/pkg/errors"
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

	mux.HandleFunc("/invalid", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>Test Page</title>
</head>
<body>
<a class="target" href="http://x.com/podcasts-test1-first/idd">invalid</a>
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
		NewGenre(1, "http://x.com/podcasts-test1-first/id1", "link #1"),
		NewGenre(2, "http://x.com/podcasts-test1-second/id2", "link #2"),
		NewGenre(3, "http://x.com/podcasts-test2-first/id3", "link #3"),
	}
}

func TestGetRequestOptions(t *testing.T) {

	options := GetRequestOptions()
	assert.NotEmpty(t, options.LookupURL)
	assert.NotEmpty(t, options.Pattern)
}

func TestGetGenres(t *testing.T) {

	ts := newTestServer()
	defer ts.Close()

	genres, _ := GetGenres(&crawler.ScraperOptions{
		LookupURL: []string{ts.URL},
		Pattern:   ".target",
	})
	mocked := getMockedGenres()

	assert.Equal(t, len(mocked), len(genres))
	for _, genre := range mocked {
		assert.Contains(t, genres, genre)
	}

	_, err := GetGenres(&crawler.ScraperOptions{
		LookupURL: []string{ts.URL + "/invalid"},
		Pattern:   ".target",
	})
	assert.Equal(t, "strconv.Atoi: parsing \"d\": invalid syntax", errors.Cause(err[0]).Error())

	_, err = GetGenres(&crawler.ScraperOptions{
		LookupURL: []string{ts.URL + "/404"},
		Pattern:   ".target",
	})
	assert.Equal(t, "Not Found", errors.Cause(err[0]).Error())
}

func TestGetGenresFromFile(t *testing.T) {

	gen, err := GetGenresFromFile("/get/invalid/path")
	assert.NotNil(t, err)
	assert.Empty(t, gen)

	func() {
		gen = []*Genre{}
		for i := 1; i <= 5; i++ {
			gen = append(gen, NewGenre(i, "http://x.com", "X"))
		}
		json, _ := json.Marshal(gen)
		ioutil.WriteFile("/tmp/genre.test.json", json, 0644)
	}()

	gen, err = GetGenresFromFile("/tmp/genre.test.json")
	assert.Nil(t, err)
	assert.Len(t, gen, 5)

	os.Remove("/tmp/genre.test.json")
}
