package show

import (
	"encoding/xml"
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

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		XMLName       xml.Name `xml:"channel"`
		Description   string   `xml:"description"`
		LastBuildDate string   `xml:"lastBuildDate"`
		Item          struct {
			Title       string `xml:"title"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func GetFeed(shows []*ShowDetails) ([]*Feed, []error) {

	urlToID := getShowsURLToID(shows)

	entities, errs := crawler.RequestEntities(
		getRequestOptions(shows),
		rssDecoder,
	)

	feeds := []*Feed{}
	for _, entity := range entities {
		feed, err := getFeedData(entity, urlToID)
		if err != nil {
			errs = append(errs, err)
		} else {
			feeds = append(feeds, feed)
		}
	}

	return feeds, errs
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

func getRequestOptions(shows []*ShowDetails) *crawler.RequestOptions {

	urls := []string{}
	for _, details := range shows {
		if details.RSS != "" {
			urls = append(urls, details.RSS)
		}
	}

	return &crawler.RequestOptions{LookupURL: urls}
}

func rssDecoder(url string, body []byte) (interface{}, error) {
	var rss RSS
	err := xml.Unmarshal(body, &rss)

	if err != nil {
		return &RSS{}, err
	}

	return rss, nil
}

func getFeedData(entity interface{}, urlToID map[string]int) (*Feed, error) {

	rss, ok := entity.(RSS)
	if !ok {
		return &Feed{}, errors.New("Invalid entity detected")
	}

	return &Feed{
		ID:          0,
		Description: rss.Channel.Description,
		LastPodcast: Podcast{
			Title:       rss.Channel.Item.Title,
			Description: rss.Channel.Item.Description,
			Published:   rss.Channel.LastBuildDate,
		},
	}, nil
}
