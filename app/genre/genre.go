package genre

import (
	// "fmt"
	// "errors"
	"encoding/json"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
	"github.com/zhikiri/uaitunes-podcasts/app/static"
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

func Save(path string, genres []*Genre) error {

	return static.Save(path, func() ([]byte, error) {

		return json.Marshal(genres)
	})
}

func GetGenresFromFile(path string) ([]*Genre, error) {

	genres := []*Genre{}

	err := static.Load(path, func(body []byte) error {

		return json.Unmarshal(body, &genres)
	})

	if err != nil {
		return []*Genre{}, err
	}

	return genres, nil
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
