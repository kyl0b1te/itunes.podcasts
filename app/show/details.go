package show

import (
	"encoding/json"
	"fmt"
	"errors"
	"io/ioutil"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
)

type ShowDetails struct {
	lookupDetails
}

type lookupDetails struct {
	ID          int
	RSS         string
	Name        string
	Genres      []string
	Image       ShowImage
	Artist      string
}

type rssDetails struct {
	Description string
	LastPodcast Podcast
}

type ShowImage struct {
	Big    string
	Small  string
	Medium string
}

type Podcast struct {
	Title       string
	Published   string
	Description string
}

type lookupResponse struct {
	Results []struct {
		CollectionId   int      `json:"collectionId"`
		ArtistName     string   `json:"artistName"`
		CollectionName string   `json:"collectionName"`
		GenreIds       []string `json:"genreIds"`
		ArtworkURL30   string   `json:"artworkURL30"`
		ArtworkURL60   string   `json:"artworkURL60"`
		ArtworkURL100  string   `json:"artworkURL100"`
		FeedURL        string   `json:"feedUrl"`
	} `json:"results"`
}

func GetDetailsRequestOptions(shows []*Show) *crawler.RequestOptions {

	urls := []string{}
	for _, show := range shows {
		url := fmt.Sprintf("%s=%d", "https://itunes.apple.com/lookup?id", show.ID)
		urls = append(urls, url)
	}

	return &crawler.RequestOptions{LookupURL: urls}
}

func GetDetails(opt *crawler.RequestOptions) ([]*ShowDetails, []error) {

	entities, errs := crawler.RequestEntities(opt, lookupDecoder)

	details := []*ShowDetails{}
	for _, entity := range entities {
		locDetails, err := getLookupDetails(entity)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		details = append(details, &ShowDetails{lookupDetails: locDetails})
	}

	return details, errs
}

func SaveDetails(file string, shows []*ShowDetails) error {

	json, err := json.MarshalIndent(shows, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, json, 0644)
}

func lookupDecoder(body []byte) (interface{}, error) {

	var res lookupResponse
	err := json.Unmarshal(body, &res)
	if err != nil {
		return &lookupResponse{}, err
	}

	return res, err
}

func getLookupDetails(entity interface{}) (lookupDetails, error) {

	res, ok := entity.(lookupResponse)
	if !ok {
		return lookupDetails{}, errors.New("Invalid entity detected")
	}

	if (len(res.Results) == 0) {
		return lookupDetails{}, errors.New("Show is not found")
	}

	apiRes := res.Results[0]

	return lookupDetails{
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
	}, nil
}
