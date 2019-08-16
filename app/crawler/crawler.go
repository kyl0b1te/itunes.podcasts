package crawler

import (
	"strings"
	"strconv"

	"github.com/gocolly/colly"
)

type RequestOptions struct {
	LookupURL string
	Pattern string
}

func GetEntities(options *RequestOptions) (map[string]string, error) {

	var err error
	entities := map[string]string{}

	col := colly.NewCollector()
	col.OnHTML(options.Pattern, func(el *colly.HTMLElement) {
		entities[el.Text] = el.Attr("href")
	})

	col.OnError(func(res *colly.Response, colErr error) {
		err = colErr
	})
	col.Visit(options.LookupURL)

	return entities, err
}

func GetEntityIDFromURL(url string) (int, error) {

	parts := strings.Split(url, "/")
	last := parts[len(parts)-1]

	return strconv.Atoi(strings.TrimPrefix(last, "id"))
}
