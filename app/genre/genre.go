package genre

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const genreURL = "https://podcasts.apple.com/us/genre/podcasts-%s/id%d"

type Genre struct {
	ID   int
	Name string
	URL string
}

type AllGenresRequestOptions struct {
	LookupURL string
	Pattern   string
}

func NewGenre(id int, name, url string) *Genre {

	return &Genre{id, name, url}
}

func GetAllGenresRequestOptions() *AllGenresRequestOptions {

	return &AllGenresRequestOptions{
		LookupURL: "https://podcasts.apple.com/us/genre/podcasts/id26",
		Pattern:   ".top-level-genre, .top-level-subgenres a[href]",
	}
}

func GetAllGenresFromWeb(options *AllGenresRequestOptions) ([]*Genre, error) {

	var err error
	genres := []*Genre{}

	collector := colly.NewCollector()
	collector.OnHTML(options.Pattern, func(element *colly.HTMLElement) {
		url := element.Attr("href")
		id := getGenreIDFromURL(url)

		genres = append(genres, NewGenre(id, element.Text, url))
	})

	collector.OnError(func(response *colly.Response, colErr error) {
		err = colErr
	})

	collector.Visit(options.LookupURL)
	collector.Wait()

	return genres, err
}

func getGenreIDFromURL(url string) int {

	src := strings.Split(url, "/")
	id, err := strconv.Atoi(strings.TrimPrefix(src[len(src)-1], "id"))
	if err != nil {
		panic(err)
	}

	return id
}
