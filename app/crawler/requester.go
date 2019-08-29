package crawler

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
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

	numb := len(opt.LookupURL)

	var wg sync.WaitGroup
	results := make(chan *RequestResult, numb)

	wg.Add(numb)
	for _, url := range opt.LookupURL {

		go func(url string) {

			results <- getEntitiesFromRequest(url, decoder)
			wg.Done()
		}(url)
	}
	wg.Wait()

	close(results)
	return results
}

func RequestEntitiesWithLimiter(opt *LimitedRequestOptions, decoder RequestDecoder) chan *RequestResult {

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

	return out
}

func getEntitiesFromRequest(url string, decoder RequestDecoder) *RequestResult {

	resp, err := http.Get(url)
	if err != nil {
		return &RequestResult{url, nil, err}
	}
	if resp.StatusCode != 200 {
		return &RequestResult{url, nil, errors.Errorf("Unreachable URL: %s", url)}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &RequestResult{url, nil, err}
	}

	res, err := decoder(url, body)
	return &RequestResult{url, res, err}
}
