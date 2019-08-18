package show

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
	"github.com/zhikiri/uaitunes-podcasts/app/genre"

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
<a class="target" href="http://x.com/sh/1">Sh1</a>
<a class="target" href="http://x.com/sh/2">Sh2</a>
<a class="target" href="http://x.com/sh/3">Sh3</a>
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

func getMockedShows() []*Show {

	return []*Show{
		NewShow(1, "http://x.com/sh/1", "Sh1"),
		NewShow(2, "http://x.com/sh/2", "Sh2"),
		NewShow(3, "http://x.com/sh/3", "Sh3"),
	}
}

func TestGetRequestOptions(t *testing.T) {

	genres := []*genre.Genre{
		genre.NewGenre(1, "http://x.com./gr/1", "Gr1"),
		genre.NewGenre(2, "http://x.com./gr/2", "Gr2"),
		genre.NewGenre(3, "http://x.com./gr/3", "Gr3"),
	}
	opt := GetRequestOptions(genres)

	for _, gr := range genres {
		assert.Contains(t, opt.LookupURL, gr.URL)
	}
	assert.NotEmpty(t, opt.Pattern)
}

func TestGetShows(t *testing.T) {

	ts := newTestServer()
	defer ts.Close()

	shows, _ := GetShows(&crawler.ScraperOptions{
		LookupURL: []string{ts.URL},
		Pattern:   ".target",
	})
	mocked := getMockedShows()

	assert.Equal(t, len(mocked), len(shows))
	for _, mockedShow := range mocked {
		assert.Contains(t, shows, mockedShow)
	}

	_, err := GetShows(&crawler.ScraperOptions{
		LookupURL: []string{ts.URL + "/invalid"},
		Pattern:   ".target",
	})
	assert.Equal(t, "strconv.Atoi: parsing \"d\": invalid syntax", errors.Cause(err[0]).Error())

	_, err = GetShows(&crawler.ScraperOptions{
		LookupURL: []string{ts.URL + "/404"},
		Pattern:   ".target",
	})
	assert.Equal(t, "Not Found", err[0].Error())
}
