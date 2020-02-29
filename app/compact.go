package main

import (
	"encoding/json"
	"strconv"

	"github.com/zhikiri/itunes.podcasts/app/genre"
	"github.com/zhikiri/itunes.podcasts/app/show"
	"github.com/zhikiri/itunes.podcasts/app/static"
)

type feedsMap = map[int]*show.Feed
type detailsMap = map[int]*show.ShowDetails
type genresMap = map[int]*genre.Genre

// CompactShow represents compacted version of the show
type CompactShow struct {
	ID       int              `json:"id"`
	ShowURL  string           `json:"show_url"`
	FeedURL  string           `json:"feed_url"`
	Name     string           `json:"name"`
	Desc     string           `json:"description"`
	Artist   string           `json:"author"`
	Language string           `json:"language"`
	Genres   []string         `json:"genres"`
	Image    CompactShowImage `json:"image"`
}

// CompactShowImage represents compact version show image
type CompactShowImage struct {
	Big    string `json:"xl"`
	Small  string `json:"xs"`
	Medium string `json:"md"`
}

func getGenresMap(genres []*genre.Genre) genresMap {
	res := make(map[int]*genre.Genre, len(genres))
	for _, genre := range genres {
		res[genre.ID] = genre
	}
	return res
}

func getDetailsMap(details []*show.ShowDetails) detailsMap {
	res := make(map[int]*show.ShowDetails, len(details))
	for _, details := range details {
		res[details.ID] = details
	}
	return res
}

func getFeedsMap(feeds []*show.Feed) feedsMap {
	res := make(map[int]*show.Feed, len(feeds))
	for _, feed := range feeds {
		res[feed.ID] = feed
	}
	return res
}

// NewCompactShow creates new instance of the compact show
func NewCompactShow(show *show.Show) *CompactShow {
	return &CompactShow{
		ID:      show.ID,
		Name:    show.Name,
		ShowURL: show.URL,
	}
}

// SetFromFeed set compact information from the show feed
func (c *CompactShow) SetFromFeed(list feedsMap) bool {
	feed, exist := list[c.ID]
	if !exist {
		return false
	}
	c.Language = feed.Language
	c.Desc = feed.Description
	return true
}

// SetFromDetails set compact information from the details along with genres
func (c *CompactShow) SetFromDetails(list detailsMap, gen genresMap) bool {
	details, exist := list[c.ID]
	if !exist {
		return false
	}
	c.Name = details.Name
	c.FeedURL = details.RSS
	c.Image.Big = details.Image.Big
	c.Image.Medium = details.Image.Medium
	c.Image.Small = details.Image.Small

	c.Genres = make([]string, 0, len(details.Genres))
	for _, idStr := range details.Genres {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		genre, ok := gen[id]
		if !ok {
			continue
		}
		c.Genres = append(c.Genres, genre.Name)
	}
	return true
}

// SaveCompactShows saves the compact show information into the file
func SaveCompactShows(path string, shows []*CompactShow) error {
	return static.Save(path, func() ([]byte, error) {
		return json.Marshal(shows)
	})
}
