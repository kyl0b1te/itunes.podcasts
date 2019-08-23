package show

import (
	"encoding/json"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
	"github.com/zhikiri/uaitunes-podcasts/app/genre"
	"github.com/zhikiri/uaitunes-podcasts/app/static"
)

type Show struct {
	ID   int
	URL  string
	Name string
}

func NewShow(id int, url string, name string) *Show {

	return &Show{id, url, name}
}

func GetShowsRequestOptions(genres []*genre.Genre) *crawler.ScraperOptions {

	urls := []string{}
	for _, genre := range genres {
		urls = append(urls, genre.URL)
	}

	return crawler.GetScraperOptions(
		urls,
		"div[id=selectedcontent] .column a[href]",
	)
}

func Save(path string, shows []*Show) error {

	return static.Save(path, func() ([]byte, error) {

		return json.Marshal(shows)
	})
}

func GetShowsFromFile(path string) ([]*Show, error) {

	shows := []*Show{}

	err := static.Load(path, func(body []byte) error {

		return json.Unmarshal(body, &shows)
	})

	if err != nil {
		return []*Show{}, err
	}

	return shows, nil
}

func GetShows(opt *crawler.ScraperOptions) ([]*Show, []error) {

	res, err := crawler.ScrapeEntities(opt)
	if len(err) > 0 {
		return []*Show{}, err
	}

	shows := []*Show{}
	for name, url := range res {

		id, err := crawler.GetEntityIDFromURL(url)
		if err != nil {
			return shows, []error{err}
		}
		shows = append(shows, NewShow(id, url, name))
	}

	return shows, []error{}
}
