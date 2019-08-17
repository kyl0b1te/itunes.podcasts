package genre

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
)

type Genre struct {
	ID   int
	Name string
	URL  string
}

func NewGenre(id int, name string, url string) *Genre {

	return &Genre{id, name, url}
}

func GenresRequestOptions() *crawler.RequestOptions {

	return &crawler.RequestOptions{
		LookupURL: "https://podcasts.apple.com/us/genre/podcasts/id26",
		Pattern:   ".top-level-genre, .top-level-subgenres a[href]",
	}
}

func SaveGenres(file string, genres []*Genre) error {

	json, err := json.Marshal(genres)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, json, 0644)
}

func GetGenres(options *crawler.RequestOptions) ([]*Genre, error) {

	genres := []*Genre{}

	entities, err := crawler.GetEntities(options)
	if err != nil {
		return genres, errors.Wrapf(err, "genres cannot be loaded from URL: %s", options.LookupURL)
	}

	for name, url := range entities {

		id, err := crawler.GetEntityIDFromURL(url)
		if err != nil {
			return genres, err
		}
		genres = append(genres, NewGenre(id, name, url))
	}

	return genres, nil
}
