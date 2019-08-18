package genre

import (
	"encoding/json"
	"io/ioutil"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
)

type Genre struct {
	ID   int
	URL  string
	Name string
}

func NewGenre(id int, url string, name string) *Genre {

	return &Genre{id, url, name}
}

func GetRequestOptions() *crawler.ScraperOptions {

	// change the country in the request to parse country specific top
	return crawler.GetScraperOptions(
		[]string{"https://podcasts.apple.com/ua/genre/podcasts/id26"},
		".top-level-genre, .top-level-subgenres a[href]",
	)
}

func SaveGenres(file string, genres []*Genre) error {

	json, err := json.Marshal(genres)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, json, 0644)
}

func GetGenres(opt *crawler.ScraperOptions) ([]*Genre, []error) {

	res, err := crawler.ScrapeEntities(opt)
	if len(err) > 0 {
		return []*Genre{}, err
	}

	genres := []*Genre{}
	for name, url := range res {

		id, err := crawler.GetEntityIDFromURL(url)
		if err != nil {
			return genres, []error{err}
		}
		genres = append(genres, NewGenre(id, url, name))
	}

	return genres, []error{}
}
