package genre

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const genreURL = "https://podcasts.apple.com/us/genre/podcasts-%s/id%d"

type Genre struct {
	ID   int
	Name string
}

type AllGenresRequestOptions struct {
	LookupURL string
	Pattern   string
}

func NewGenre(id int, name string) *Genre {

	return &Genre{id, name}
}

func newGenreByURL(url string) *Genre {

	src := strings.Split(url, "/")
	name := strings.TrimPrefix(src[len(src)-2], "podcasts-")

	id, err := strconv.Atoi(strings.TrimPrefix(src[len(src)-1], "id"))
	if err != nil {
		panic(err)
	}

	return NewGenre(id, name)
}

func (g Genre) GetURL() string {

	return fmt.Sprintf(genreURL, strings.ToLower(g.Name), g.ID)
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
		link := element.Attr("href")
		genres = append(genres, newGenreByURL(link))
	})

	collector.OnError(func(response *colly.Response, colErr error) {
		err = colErr
	})

	collector.Visit(options.LookupURL)
	collector.Wait()

	return genres, err
}
