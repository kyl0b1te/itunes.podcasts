package crawler

import (
	"strconv"
	"strings"
	"sync"
	// "fmt"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

type ScraperOptions struct {
	LookupURL []string
	Pattern   string
}

type ScrapeResult struct {
	Entities map[string]string
	Errors   []error
}

func GetScraperOptions(url []string, pattern string) *ScraperOptions {

	return &ScraperOptions{url, pattern}
}

func ScrapeEntities(opt *ScraperOptions) (map[string]string, []error) {

	var wg sync.WaitGroup
	resCh := make(chan *ScrapeResult, len(opt.LookupURL))

	wg.Add(len(opt.LookupURL))
	for _, url := range opt.LookupURL {

		go func(url string) {

			resCh <- getEntitiesFromHTML(url, opt.Pattern)
			wg.Done()
		}(url)
	}
	wg.Wait()

	close(resCh)

	res := map[string]string{}
	err := []error{}

	for scrape := range resCh {

		if len(scrape.Errors) > 0 {
			err = append(err, scrape.Errors...)
		}

		for name, url := range scrape.Entities {
			if _, ok := res[name]; !ok {
				res[name] = url
			}
		}
	}

	return res, err
}

func getEntitiesFromHTML(url string, pattern string) *ScrapeResult {

	errs := []error{}
	res := map[string]string{}

	col := colly.NewCollector()
	col.OnHTML(pattern, func(el *colly.HTMLElement) {
		res[el.Text] = el.Attr("href")
	})

	col.OnError(func(resp *colly.Response, err error) {
		// todo : wrap error
		errs = append(errs, err)
	})
	col.Visit(url)

	return &ScrapeResult{res, errs}
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
