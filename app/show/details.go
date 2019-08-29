package show

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/zhikiri/uaitunes-podcasts/app/crawler"
	"github.com/zhikiri/uaitunes-podcasts/app/static"
)

type ShowDetails struct {
	ID          int
	RSS         string
	Name        string
	Genres      []string
	Artist      string
	Image       ShowImage
	Description string
	LastPodcast Podcast
}

type ShowImage struct {
	Big    string
	Small  string
	Medium string
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

func GetDetailsRequestOptions(shows []*Show) *crawler.LimitedRequestOptions {

	urls := []string{}
	for _, show := range shows {
		url := fmt.Sprintf("%s=%d", "https://itunes.apple.com/lookup?id", show.ID)
		urls = append(urls, url)
	}

	return &crawler.LimitedRequestOptions{
		LookupURL: urls,
		Duration:  5 * time.Second,
	}
}

func GetDetails(opt *crawler.LimitedRequestOptions) ([]*ShowDetails, []error) {

	details := []*ShowDetails{}
	errs := []error{}

	out := crawler.RequestEntitiesWithLimiter(opt, lookupDecoder)
	for en := range out {
		if en.Error != nil {
			errs = append(errs, en.Error)
			continue
		}
		det, err := getLookupDetails(en)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		details = append(details, det)
	}

	return details, errs
}

func SaveDetails(path string, details []*ShowDetails) error {

	return static.Save(path, func() ([]byte, error) {

		return json.Marshal(details)
	})
}

func GetShowDetailsFromFile(path string) ([]*ShowDetails, error) {

	details := []*ShowDetails{}

	err := static.Load(path, func(body []byte) error {

		return json.Unmarshal(body, &details)
	})

	if err != nil {
		return []*ShowDetails{}, err
	}

	return details, nil
}

func lookupDecoder(url string, body []byte) (interface{}, error) {

	var res lookupResponse
	err := json.Unmarshal(body, &res)
	if err != nil {
		return &lookupResponse{}, err
	}

	return res, err
}

func getLookupDetails(entity interface{}) (*ShowDetails, error) {

	res, ok := entity.(lookupResponse)
	if !ok {
		return &ShowDetails{}, errors.New("Invalid entity detected")
	}

	if len(res.Results) == 0 {
		return &ShowDetails{}, errors.New("Show is not found")
	}

	apiRes := res.Results[0]

	return &ShowDetails{
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
		Description: "",
		LastPodcast: Podcast{
			Title:       "",
			Description: "",
			Published:   "",
		},
	}, nil
}
