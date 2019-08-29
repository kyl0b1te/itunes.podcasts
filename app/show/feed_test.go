package show

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const feedTemplate = `
<rss>
	<channel>
		<description>desc_X</description>
		<language>en</language>
		<lastBuildDate>pub_X</lastBuildDate>
		<item>
			<title>item_title_X</title>
			<description>item_desc_X</description>
		</item>
	</channel>
</rss>
`

func newFeedTestServer() *httptest.Server {

	mux := http.NewServeMux()
	for url, data := range getTestFeedEndpoints() {
		func(url string, data []byte) {
			mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/xml")
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

	mux.HandleFunc("/invalid", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(200)
		w.Write([]byte("<invalid>Data</invalid>"))
	})

	return httptest.NewServer(mux)
}

func getTestFeedEndpoints() map[string][]byte {

	return map[string][]byte{
		"/show/1": []byte(strings.ReplaceAll(feedTemplate, "X", "1")),
		"/show/2": []byte(strings.ReplaceAll(feedTemplate, "X", "2")),
		"/show/3": []byte(strings.ReplaceAll(feedTemplate, "X", "3")),
	}
}

func getTestShowDetails(ts *httptest.Server) []*ShowDetails {

	details := []*ShowDetails{}
	for i := 1; i <= 3; i++ {
		det := &ShowDetails{ID: i, RSS: fmt.Sprintf("%s/show/%d", ts.URL, i)}
		details = append(details, det)
	}
	return details
}

func TestGetFeed(t *testing.T) {

	ts := newFeedTestServer()
	defer ts.Close()

	details := getTestShowDetails(ts)
	list, errs := GetFeed(details)

	assert.Empty(t, errs)
	for _, feed := range list {
		assert.Contains(t, []int{1, 2, 3}, feed.ID)
		assert.Equal(t, fmt.Sprintf("desc_%d", feed.ID), feed.Description)
		assert.Equal(t, fmt.Sprintf("pub_%d", feed.ID), feed.LastPodcast.Published)
		assert.Equal(t, fmt.Sprintf("item_title_%d", feed.ID), feed.LastPodcast.Title)
		assert.Equal(t, fmt.Sprintf("item_desc_%d", feed.ID), feed.LastPodcast.Description)
		assert.Equal(t, "en", feed.Language)
	}

	details = []*ShowDetails{
		&ShowDetails{ID: 1, RSS: ts.URL + "/404"},
	}
	_, errs = GetFeed(details)
	assert.NotEmpty(t, errs)
	msg := fmt.Sprintf("Unreachable URL: %s/404", ts.URL)
	assert.Equal(t, msg, errors.Cause(errs[0]).Error())

	details = []*ShowDetails{
		&ShowDetails{ID: 1, RSS: ts.URL + "/invalid"},
	}
	_, errs = GetFeed(details)
	assert.NotEmpty(t, errs)
	msg = "expected element type <rss> but have <invalid>"
	assert.Equal(t, msg, errors.Cause(errs[0]).Error())
}
