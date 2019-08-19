package crawler

import (
	"net/http"
	"io/ioutil"
	"sync"
	// "fmt"
)

type RequestResult struct {
	Entity interface{}
	Error error
}

type RequestOptions struct {
	LookupURL []string
}

type RequestDecoder func(body []byte) (interface{}, error)

func RequestEntities(opt *RequestOptions, decoder RequestDecoder) ([]interface{}, []error) {

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

	res := make([]interface{}, 0, urlNumber)
	err := []error{}

	for result := range resCh {

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
		return &RequestResult{nil, err}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &RequestResult{nil, err}
	}

	res, err := decoder(body)
	return &RequestResult{res, err}
}
