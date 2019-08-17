package show

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
	"github.com/zhikiri/uaitunes-podcasts/app/genre"
)

type Show struct {
	ID     int
	Name   string
	Artist string
	RSS    string
	Genres []string
	Image  ShowImage
}

type ShowImage struct {
	Small  string
	Medium string
	Big    string
}

type ShowRequestOptions struct {
	crawler.RequestOptions
	ShowDetailsURL string
}

type showDetailsResponse struct {
	Results []struct {
		CollectionId   int    `json:"collectionId"`
		ArtistId       int    `json:"artistId"`
		ArtistName     string `json:"artistName"`
		CollectionName string `json:"collectionName"`
		GenreIds       []string  `json:"genreIds"`
		ArtworkURL30   string `json:"artworkURL30"`
		ArtworkURL60   string `json:"artworkURL60"`
		ArtworkURL100  string `json:"artworkURL100"`
		FeedURL        string `json:"feedUrl"`
	} `json:"results"`
}

func ShowsRequestOptions(genre *genre.Genre) *ShowRequestOptions {

	return &ShowRequestOptions{
		RequestOptions: crawler.RequestOptions{
			LookupURL: genre.URL,
			Pattern:   "div[id=selectedcontent] .column a[href]",
		},
		ShowDetailsURL: "https://itunes.apple.com/lookup?id=",
	}
}

func GetShows(options *ShowRequestOptions) ([]*Show, []error) {

	entities, err := crawler.GetEntities(&options.RequestOptions)
	if err != nil {
		return []*Show{}, []error{err}
	}

	resCh, errCh := getShowsFromEntities(entities, options)

	shows := []*Show{}
	for show := range resCh {
		shows = append(shows, show)
	}

	errors := []error{}
	for err := range errCh {
		errors = append(errors, err)
	}

	return shows, errors
}

func SaveShows(file string, shows []*Show) error {

	json, err := json.Marshal(shows)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, json, 0644)
}

func getShowsFromEntities(entities map[string]string, options *ShowRequestOptions) (chan *Show, chan error) {

	entLen := len(entities)
	errCh := make(chan error, entLen)
	resCh := make(chan *Show, entLen)

	var wg sync.WaitGroup
	wg.Add(entLen)
	for _, url := range entities {

		go func(url string) {

			if show, err := getShowDetails(url, options); err != nil {
				errCh <- err
			} else {
				resCh <- show
			}
			wg.Done()
		}(url)
	}
	wg.Wait()
	close(errCh)
	close(resCh)

	return resCh, errCh
}

func getShowDetails(showURL string, options *ShowRequestOptions) (*Show, error) {

	id, err := crawler.GetEntityIDFromURL(showURL)
	if err != nil {
		return &Show{}, err
	}

	url := fmt.Sprintf("%s%d", options.ShowDetailsURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return &Show{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Show{}, err
	}

	var details showDetailsResponse
	err = json.Unmarshal(body, &details)
	if err != nil {
		return &Show{}, err
	}

	return getShowFromResponse(details), nil
}

func getShowFromResponse(details showDetailsResponse) *Show {

	res := details.Results[0]
	return &Show{
		ID:     res.CollectionId,
		Name:   res.CollectionName,
		Artist: res.ArtistName,
		RSS:    res.FeedURL,
		Genres: res.GenreIds,
		Image: ShowImage{
			Small:  res.ArtworkURL30,
			Medium: res.ArtworkURL60,
			Big:    res.ArtworkURL100,
		},
	}
}
