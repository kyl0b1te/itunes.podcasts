package show

import (
	"fmt"
	"encoding/json"
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
		GenreIds       []int  `json:"genreIds"`
		ArtworkURL30   string `json:"artworkURL30"`
		ArtworkURL60   string `json:"artworkURL60"`
		ArtworkURL100  string `json:"artworkURL100"`
	} `json:"results"`
}

func NewShow(id int, name string, artist string) *Show {

	return &Show{id, name, artist}
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

func GetShows(options *ShowRequestOptions) ([]*Show, error) {

	shows := []*Show{}

	entities, err := crawler.GetEntities(&options.RequestOptions)
	if err != nil {
		return shows, err
	}

	ch := make(chan *Show, len(entities))
	var wg sync.WaitGroup

	for _, url := range entities {

		id, err := crawler.GetEntityIDFromURL(url)
		if err != nil {
			return shows, err
		}

		wg.Add(1)
		go func() {
			show, err := getShowDetails(id, options)
			if err != nil {
				// todo : store errors in a channel
				fmt.Println(err)
			}
			ch <- show
			wg.Done()
		}()
	}
	wg.Wait()
	close(ch)

	for show := range ch {
		shows = append(shows, show)
	}

	return shows, nil
}

func getShowDetails(id int, options *ShowRequestOptions) (*Show, error) {

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

	return &Show{
		ID:     details.Results[0].CollectionId,
		Name:   details.Results[0].CollectionName,
		Artist: details.Results[0].ArtistName,
	}, nil
}
