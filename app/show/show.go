package show

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
	"github.com/zhikiri/uaitunes-podcasts/app/genre"
)

type Show struct {
	ID          int
	Name        string
	Artist      string
	RSS         string
	Genres      []string
	Image       ShowImage
	Description string
	LastPodcast Podcast
	Language	  string
}

type ShowImage struct {
	Small  string
	Medium string
	Big    string
}

type Podcast struct {
	Title       string
	Description string
}

type ShowRequestOptions struct {
	crawler.RequestOptions
	ShowDetailsURL string
}

type apiResponse struct {
	Results []struct {
		CollectionId   int      `json:"collectionId"`
		ArtistId       int      `json:"artistId"`
		ArtistName     string   `json:"artistName"`
		CollectionName string   `json:"collectionName"`
		GenreIds       []string `json:"genreIds"`
		ArtworkURL30   string   `json:"artworkURL30"`
		ArtworkURL60   string   `json:"artworkURL60"`
		ArtworkURL100  string   `json:"artworkURL100"`
		FeedURL        string   `json:"feedUrl"`
	} `json:"results"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		XMLName     xml.Name `xml:"channel"`
		Description string   `xml:"description"`
		Item        struct {
			Title       string   `xml:"title"`
			Description string   `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
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

	json, err := json.MarshalIndent(shows, "", "\t")
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

	details, err := getShowFromAPI(id, options)
	if err != nil {
		return &Show{}, err
	}

	rss := &RSS{}
	if details.Results[0].FeedURL != "" {
		rss, _ = getShowFromRSS(details.Results[0].FeedURL)
	}

	return newShow(details, rss), nil
}

func getShowFromAPI(id int, options *ShowRequestOptions) (*apiResponse, error) {

	url := fmt.Sprintf("%s%d", options.ShowDetailsURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return &apiResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &apiResponse{}, err
	}

	var details apiResponse
	err = json.Unmarshal(body, &details)
	if err != nil {
		return &apiResponse{}, err
	}

	return &details, nil
}

func getShowFromRSS(url string) (*RSS, error) {

	resp, err := http.Get(url)
	if err != nil {
		return &RSS{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &RSS{}, err
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return &RSS{}, err
	}

	return &rss, nil
}

func newShow(api *apiResponse, rss *RSS) *Show {

	apiRes := api.Results[0]

	return &Show{
		ID:     apiRes.CollectionId,
		Name:   apiRes.CollectionName,
		Artist: apiRes.ArtistName,
		RSS:    apiRes.FeedURL,
		Genres: apiRes.GenreIds,
		Image: ShowImage{
			Small:  apiRes.ArtworkURL30,
			Medium: apiRes.ArtworkURL60,
			Big:    apiRes.ArtworkURL100,
		},
		Description: rss.Channel.Description,
		LastPodcast: Podcast{
			Title:       rss.Channel.Item.Title,
			Description: rss.Channel.Item.Description,
		},
		Language: "",
	}
}
