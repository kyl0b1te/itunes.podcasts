package crawler

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type RequestResult struct {
	URL    string
	Entity interface{}
	Error  error
}

type RequestOptions struct {
	LookupURL []string
}

type LimitedRequestOptions struct {
	LookupURL []string
	Duration  time.Duration
}

type RequestDecoder func(url string, body []byte) (interface{}, error)

func RequestEntities(opt *RequestOptions, decoder RequestDecoder) chan *RequestResult {

	urlNumber := len(opt.LookupURL)

	var wg sync.WaitGroup
	resCh := make(chan *RequestResult, urlNumber)

	wg.Add(urlNumber)
	for _, url := range opt.LookupURL {

		go func(url string) {

			resCh <- getEntitiesFromRequest(url, decoder)
			wg.Done()
		}(url)
	}
	wg.Wait()

	close(resCh)

	return resCh
}

func RequestEntitiesWithLimiter(opt *LimitedRequestOptions, decoder RequestDecoder) ([]interface{}, []error) {

	urls := len(opt.LookupURL)

	in := make(chan string, urls)
	out := make(chan *RequestResult, urls)

	for _, url := range opt.LookupURL {
		in <- url
	}
	close(in)

	limiter := time.Tick(opt.Duration)

	go func(in chan string, out chan *RequestResult) {
		i := 1
		for url := range in {
			<-limiter
			log.Printf("Requesting (%d/%d) - %s", i, urls, url)
			out <- getEntitiesFromRequest(url, decoder)
			i++
		}
		close(out)
	}(in, out)

	res := make([]interface{}, 0, urls)
	err := []error{}

	for result := range out {

		if result.Error != nil {
			err = append(err, result.Error)
		}

		if result.Entity != nil {
			res = append(res, result.Entity)
		}
	}

	return res, err
}

func getEntitiesFromRequest(url string, decoder RequestDecoder) *RequestResult {

	resp, err := http.Get(url)
	if err != nil {
		return &RequestResult{url, nil, err}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &RequestResult{url, nil, err}
	}

	res, err := decoder(url, body)
	return &RequestResult{url, res, err}
}
