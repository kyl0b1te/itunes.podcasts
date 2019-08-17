package show

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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

	for _, show := range getMockedShows() {

		resp := getMockedShowResponse(show)

		mux.HandleFunc(
			fmt.Sprintf("/show/%d", show.ID),
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write(resp)
			},
		)
	}

	mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(404)
		w.Write([]byte("<p>error</p>"))
	})

	return httptest.NewServer(mux)
}

func getMockedShowResponse(show *Show) []byte {

	msg := fmt.Sprintf(
		`{"results": [
        {
            "collectionId": %d,
            "artistName": "%s",
						"collectionName": "%s",
						"genreIds": [],
						"artworkURL30": "",
						"artworkURL60": "",
						"artworkURL100": "",
						"feedUrl": ""
        }
		]}`,
		show.ID,
		show.Artist,
		show.Name,
	)
	return []byte(msg)
}

func getNewShow(id int, name string, artist string) *Show {

	return &Show{
		ID:     id,
		Name:   name,
		Artist: artist,
		RSS:    "",
		Genres: []string{},
		Image: ShowImage{
			Small:  "",
			Medium: "",
			Big:    "",
		},
		Description: "",
	}
}

func getMockedShows() []*Show {

	return []*Show{
		getNewShow(1, "Show #1", "Artist Show #1"),
		getNewShow(2, "Show #2", "Artist Show #2"),
		getNewShow(3, "Show #3", "Artist Show #3"),
	}
}

func TestGetShows(t *testing.T) {

	ts := newTestServer()
	defer ts.Close()

	options := &ShowRequestOptions{
		RequestOptions: crawler.RequestOptions{
			LookupURL: ts.URL,
			Pattern:   ".target",
		},
		ShowDetailsURL: ts.URL + "/show/",
	}

	shows, _ := GetShows(options)
	mocked := getMockedShows()

	assert.Equal(t, len(mocked), len(shows))
	for _, mockedShow := range mocked {
		assert.Contains(t, shows, mockedShow)
	}

	options.RequestOptions.LookupURL = ts.URL + "/404"
	_, errors := GetShows(options)
	assert.Equal(t, "Not Found", errors[0].Error())
}
