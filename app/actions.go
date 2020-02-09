package main

import (
	"fmt"
	"path"
	"time"

	"github.com/zhikiri/itunes.podcasts/app/genre"
	"github.com/zhikiri/itunes.podcasts/app/show"
)

func actionGenres(out string) {
	fmt.Println("Starting genres loading")
	genres, errs := genre.GetGenres(genre.GetRequestOptions())
	stopOnErrors(errs)

	fmt.Println("Genres loaded", len(genres))
	err := genre.Save(path.Join(out, "genres.json"), genres)
	stopOnError(err)
}

func actionShows(genrePath string, out string) {
	fmt.Println("Starting shows loading")
	genres, err := genre.GetGenresFromFile(genrePath)
	stopOnError(err)

	shows, errs := show.GetShows(show.GetShowsRequestOptions(genres))
	stopOnErrors(errs)

	fmt.Println("Shows loaded", len(shows))
	err = show.Save(path.Join(out, "shows.json"), shows)
	stopOnError(err)
}

func actionDetails(showPath string, chunk int, out string) {
	fmt.Println("Starting details loading")
	shows, err := show.GetShowsFromFile(showPath)
	stopOnError(err)
	fmt.Println("Shows total", len(shows))

	file := path.Join(out, "shows.details.json")
	cache, _ := show.GetShowDetailsFromFile(file)
	fmt.Println("Details found", len(cache))

	inCache := make(map[int]int, len(cache))
	for _, show := range cache {
		inCache[show.ID] = 1
	}

	fresh := make([]*show.Show, 0, chunk)
	for _, show := range shows {
		if _, ok := inCache[show.ID]; !ok && len(fresh) < chunk {
			fresh = append(fresh, show)
		}
	}

	details, errs := show.GetDetails(show.GetDetailsRequestOptions(fresh, 5*time.Second))
	stopOnErrors(errs)

	fmt.Println("Details loaded", len(details))
	cache = append(cache, details...)
	err = show.SaveDetails(file, cache)
	stopOnError(err)
}

func actionFeed(detailPath string, out string) {
	fmt.Println("Starting feed loading")
	details, err := show.GetShowDetailsFromFile(detailPath)
	stopOnError(err)

	fmt.Println("Details found", len(details))
	feeds, errs := show.GetFeed(details)

	fmt.Println("Feeds loaded", len(feeds))
	err = show.SaveFeed(path.Join(out, "shows.feed.json"), feeds)
	if err != nil {
		errs = append(errs, err)
	}
	stopOnErrors(errs)
}
