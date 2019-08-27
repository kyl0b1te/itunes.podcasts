package show

import (
	"errors"
	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
)

type Feed struct {
	ID          int
	Description string
	LastPodcast Podcast
}

type Podcast struct {
	Title       string
	Published   string
	Description string
}

type RSS struct {}

func GetFeed(shows []*ShowDetails) ([]*Feed, []error) {

	urlToID := getShowsURLToID(shows)

	entities, errs := crawler.RequestEntities(
		getRequestOptions(shows),
		rssDecoder(urlToID),
	)

	feeds := []*Feed{}
	for _, entity := range entities {
		// todo : parse feed here
	}
}

func getShowsURLToID(shows []*ShowDetails) map[string]int {

	res := map[string]int{}
	for _, details := range shows {
		if details.RSS != "" {
			res[details.RSS] = details.ID
		}
	}
	return res
}

func getRequestOptions(shows []*ShowDetails) *crawler.getRequestOptions {

	urls := []string{}
	for _, details := range shows {
		if details.RSS != "" {
			urls = append(urls, details.RSS)
		}
	}

	return &crawler.RequestOptions{urls}
}

func rssDecoder(urlToID map[string]int) crawler.RequestDecoder {

	return func(url string, body []byte) (interface{}, error) {
		var feed RSS
		err := json.Unmarshal(body, &feed)

		if err != nil {
			return &Feed{}, err
		}

		if id, ok := urlToID[url]; !ok {
			return &Feed{}, errors.New("Cannot resolve feed URL show id")
		} else {
			feed.ID = id
		}

		return feed, nil
	}
}
