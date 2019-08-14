package genre

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

	return httptest.NewServer(mux)
}

func getMockedGenres() []*Genre {

	return []*Genre{
		&Genre{1, "test1-first"},
		&Genre{2, "test1-second"},
		&Genre{3, "test2-first"},
	}
}

func TestGetURL(t *testing.T) {

	genre := NewGenre(1, "NAME")
	url := genre.GetURL()

	assert.Equal(t, "https://podcasts.apple.com/us/genre/podcasts-name/id1", url)
}

func TestGetGenresFromWeb(t *testing.T) {

	ts := newTestServer()
	defer ts.Close()

	actualGenres := GetAllGenresFromWeb(
		&AllGenresRequestOptions{LookupURL: ts.URL, Pattern: ".target"},
	)
	expectedGenres := getMockedGenres()

	assert.Equal(t, 3, len(actualGenres))
	for _, genre := range expectedGenres {
		assert.Contains(t, actualGenres, genre)
	}
}
