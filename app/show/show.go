package show

import (
	"encoding/json"
	"io/ioutil"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
	"github.com/zhikiri/uaitunes-podcasts/app/genre"
)

type Show struct {
	ID   int
	Name string
	URL  string
}

func NewShow(id int, name string, url string) *Show {

	return &Show{id, name, url}
}

func GetRequestOptions(genres []*genre.Genre) *crawler.ScraperOptions {

	urls := []string{}
	for _, genre := range genres {
		urls = append(urls, genre.URL)
	}

	return crawler.GetScraperOptions(
		urls,
		"div[id=selectedcontent] .column a[href]",
	)
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
		shows = append(shows, NewShow(id, name, url))
	}

	return shows, []error{}
}

func SaveShows(file string, shows []*Show) error {

	json, err := json.MarshalIndent(shows, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, json, 0644)
}
