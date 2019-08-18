package show

import (
	"encoding/json"
	"io/ioutil"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
	"github.com/zhikiri/uaitunes-podcasts/app/genre"

	"github.com/pkg/errors"
)

type Show struct {
	ID     int
	Name   string
	URL    string
	Genres []int
}

func NewShow(id int, name string, url string, genres []int) *Show {

	return &Show{id, name, url, genres}
}

func GetRequestOptions(genre *genre.Genre) *crawler.RequestOptions {

	return crawler.GetRequestOptions(genre.URL, "div[id=selectedcontent] .column a[href]")
}

func GetShows(genres []*genre.Genre) ([]*Show, error) {

	res := map[int]*Show{}
	for _, genre := range genres {

		opt := GetRequestOptions(genre)
		if err := loadShows(opt, res); err != nil {
			return []*Show{}, err
		}
	}

	shows := make([]*Show, 0, len(res))
	for _, show := range res {
		shows = append(shows, show)
	}

	return shows, nil
}

func loadShows(opt *crawler.RequestOptions, shows map[int]*Show) error {

	entities, err := crawler.GetEntities(opt)
	if err != nil {
		return errors.Wrapf(err, "shows cannot be loaded from URL: %s", opt.LookupURL)
	}

	for name, url := range entities {

		id, err := crawler.GetEntityIDFromURL(url)
		if err != nil {
			return err
		}

		if _, ok := shows[id]; !ok {
			shows[id] = NewShow(id, name, url, []int{})
		}
	}

	return nil
}

func SaveShows(file string, shows []*Show) error {

	json, err := json.MarshalIndent(shows, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, json, 0644)
}
