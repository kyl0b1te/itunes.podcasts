package crawler

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

type RequestOptions struct {
	LookupURL string
	Pattern   string
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

	id, err := strconv.Atoi(strings.TrimPrefix(last, "id"))
	if err != nil {
		return 0, errors.Wrapf(err, "ID cannot be parsed from URL: %s", url)
	}

	return id, nil
}
