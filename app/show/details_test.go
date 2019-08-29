package show

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const detailsTemplate = `
{
	"results": [{
		"collectionId": X,
		"artistName": "arts_X",
		"collectionName": "col_X",
		"genreIds": ["a", "b", "c"],
		"artworkURL30": "30_X",
		"artworkURL60": "60_X",
		"artworkURL100": "100_X",
		"feedUrl": "feed_X"
	}]
}
`

func newDetailsTestServer() *httptest.Server {

	mux := http.NewServeMux()
	for url, data := range getTestDetailsEndpoints() {
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

	mux.HandleFunc("/invalid", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")
		w.WriteHeader(200)
		w.Write([]byte("{\"data\": 1}"))
	})

	return httptest.NewServer(mux)
}

func getTestDetailsEndpoints() map[string][]byte {

	return map[string][]byte{
		"/show/1": []byte(strings.ReplaceAll(detailsTemplate, "X", "1")),
		"/show/2": []byte(strings.ReplaceAll(detailsTemplate, "X", "2")),
		"/show/3": []byte(strings.ReplaceAll(detailsTemplate, "X", "3")),
	}
}

func getTestDetailsURL(ts *httptest.Server) []string {

	urls := []string{}
	for url, _ := range getTestDetailsEndpoints() {
		urls = append(urls, ts.URL+url)
	}
	return urls
}

func TestGetDetails(t *testing.T) {

	ts := newDetailsTestServer()
	defer ts.Close()

	opt := &crawler.LimitedRequestOptions{
		LookupURL: getTestDetailsURL(ts),
		Duration:  time.Second,
	}
	list, errs := GetDetails(opt)

	assert.Empty(t, errs)
	for _, det := range list {
		assert.Contains(t, []int{1, 2, 3}, det.ID)
		assert.Equal(t, fmt.Sprintf("feed_%d", det.ID), det.RSS)
		assert.Equal(t, fmt.Sprintf("col_%d", det.ID), det.Name)
		assert.Equal(t, fmt.Sprintf("arts_%d", det.ID), det.Artist)
		assert.Equal(t, fmt.Sprintf("100_%d", det.ID), det.Image.Big)
		assert.Equal(t, fmt.Sprintf("60_%d", det.ID), det.Image.Medium)
		assert.Equal(t, fmt.Sprintf("30_%d", det.ID), det.Image.Small)
	}

	opt = &crawler.LimitedRequestOptions{
		LookupURL: []string{ts.URL + "/404"},
		Duration:  time.Second,
	}
	_, errs = GetDetails(opt)
	assert.NotEmpty(t, errs)
	msg := fmt.Sprintf("Unreachable URL: %s/404", ts.URL)
	assert.Equal(t, msg, errors.Cause(errs[0]).Error())

	opt = &crawler.LimitedRequestOptions{
		LookupURL: []string{ts.URL + "/invalid"},
		Duration:  time.Second,
	}
	_, errs = GetDetails(opt)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "Show is not found", errors.Cause(errs[0]).Error())
}

func TestGetDetailsRequestOptions(t *testing.T) {

	shows := []*Show{
		NewShow(1, "http://x.com", "1"),
		NewShow(2, "http://x.com", "2"),
		NewShow(3, "http://x.com", "3"),
	}
	opt := GetDetailsRequestOptions(shows)

	for _, sho := range shows {
		url := fmt.Sprintf("%s=%d", "https://itunes.apple.com/lookup?id", sho.ID)
		assert.Contains(t, opt.LookupURL, url)
	}
	assert.Equal(t, time.Second*5, opt.Duration)
}

func TestGetShowDetailsFromFile(t *testing.T) {

	path := "/tmp/show-details.test.json"

	det, err := GetShowDetailsFromFile("/get/invalid/path")
	assert.NotNil(t, err)
	assert.Empty(t, det)

	func() {
		det = []*ShowDetails{}
		for i := 1; i <= 5; i++ {
			det = append(det, &ShowDetails{ID: i})
		}
		json, _ := json.Marshal(det)
		ioutil.WriteFile(path, json, 0644)
	}()

	det, err = GetShowDetailsFromFile(path)
	assert.Nil(t, err)
	assert.Len(t, det, 5)

	os.Remove(path)
}
