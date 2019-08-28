package show

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
	"github.com/zhikiri/uaitunes-podcasts/app/static"
)

type Feed struct {
	ID          int
	Language    string
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
		Language      string   `xml:"language"`
		LastBuildDate string   `xml:"lastBuildDate"`
		Item          struct {
			Title       string `xml:"title"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func GetFeed(shows []*ShowDetails) ([]*Feed, []error) {

	urlToID := getShowsURLToID(shows)

	feedList := make([]*Feed, 0, len(shows))
	errs := []error{}

	out := crawler.RequestEntities(getRequestOptions(shows), rssDecoder)
	for entity := range out {
		if entity.Error != nil {
			errs = append(errs, entity.Error)
			continue
		}
		feed, err := getFeedData(entity.Entity, entity.URL, urlToID)
		if err != nil {
			errs = append(errs, err)
		} else {
			feedList = append(feedList, feed)
		}
	}

	return feedList, errs
}

func SaveFeed(path string, feed []*Feed) error {

	return static.Save(path, func() ([]byte, error) {

		return json.Marshal(feed)
	})
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

func getFeedData(entity interface{}, url string, urlToID map[string]int) (*Feed, error) {

	rss, ok := entity.(RSS)
	if !ok {
		return &Feed{}, errors.New("Invalid entity detected")
	}

	var id int
	if id, ok = urlToID[url]; !ok {
		return &Feed{}, errors.New("Cannot retrieve the show id")
	}

	lang := strings.Split(strings.ToLower(rss.Channel.Language), "-")[0]

	return &Feed{
		ID:          id,
		Description: rss.Channel.Description,
		Language:    lang,
		LastPodcast: Podcast{
			Title:       rss.Channel.Item.Title,
			Description: rss.Channel.Item.Description,
			Published:   rss.Channel.LastBuildDate,
		},
	}, nil
}
