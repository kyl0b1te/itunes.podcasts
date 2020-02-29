package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhikiri/itunes.podcasts/app/genre"
	"github.com/zhikiri/itunes.podcasts/app/show"
)

func TestGetGenresMap(t *testing.T) {
	res := getGenresMap([]*genre.Genre{
		&genre.Genre{ID: 1, Name: "test #1"},
		&genre.Genre{ID: 2, Name: "test #2"},
		&genre.Genre{ID: 3, Name: "test #3"},
	})
	assert.Len(t, res, 3)
	for id := range res {
		assert.Contains(t, []int{1, 2, 3}, id)
		assert.Equal(t,
			&genre.Genre{ID: id, Name: fmt.Sprintf("test #%d", id)},
			res[id],
		)
	}
}

func TestGetDetailsMap(t *testing.T) {
	res := getDetailsMap([]*show.ShowDetails{
		&show.ShowDetails{ID: 1, Name: "1"},
		&show.ShowDetails{ID: 2, Name: "2"},
		&show.ShowDetails{ID: 3, Name: "3"},
	})
	assert.Len(t, res, 3)
	for id := range res {
		assert.Contains(t, []int{1, 2, 3}, id)
		assert.Equal(t,
			&show.ShowDetails{ID: id, Name: fmt.Sprintf("%d", id)},
			res[id],
		)
	}
}

func TestGetFeedsMap(t *testing.T) {
	res := getFeedsMap([]*show.Feed{
		&show.Feed{ID: 1, Description: "1"},
		&show.Feed{ID: 2, Description: "2"},
		&show.Feed{ID: 3, Description: "3"},
	})
	assert.Len(t, res, 3)
	for id := range res {
		assert.Contains(t, []int{1, 2, 3}, id)
		assert.Equal(t,
			&show.Feed{ID: id, Description: fmt.Sprintf("%d", id)},
			res[id],
		)
	}
}

func TestNewCompactShow(t *testing.T) {
	show := &show.Show{ID: 1, Name: "Test", URL: "/test"}
	res := NewCompactShow(show)
	assert.Equal(t, show.ID, res.ID)
	assert.Equal(t, show.Name, res.Name)
	assert.Equal(t, show.URL, res.ShowURL)
}

func TestCompactSetFromFeed(t *testing.T) {
	com := &CompactShow{ID: 1}
	res := com.SetFromFeed(map[int]*show.Feed{
		3: &show.Feed{},
	})
	assert.False(t, res)

	src := map[int]*show.Feed{
		1: &show.Feed{Language: "test", Description: "test"},
	}
	res = com.SetFromFeed(src)
	assert.True(t, res)
	assert.Equal(t, src[1].Language, com.Language)
	assert.Equal(t, src[1].Description, com.Desc)
}

func TestCompactSetFromDetails(t *testing.T) {
	com := &CompactShow{ID: 1}
	res := com.SetFromDetails(
		map[int]*show.ShowDetails{3: &show.ShowDetails{}},
		map[int]*genre.Genre{},
	)
	assert.False(t, res)

	src := map[int]*show.ShowDetails{
		1: &show.ShowDetails{
			Name: "Name",
			RSS:  "RSS",
			Image: show.ShowImage{
				Big:    "Big",
				Medium: "Medium",
				Small:  "Small",
			},
			Genres: []string{"1", "2", "3"},
		},
	}
	res = com.SetFromDetails(src, map[int]*genre.Genre{})
	assert.True(t, res)
	assert.Equal(t, src[1].Name, com.Name)
	assert.Equal(t, src[1].RSS, com.FeedURL)
	assert.Equal(t, src[1].Image.Big, com.Image.Big)
	assert.Equal(t, src[1].Image.Medium, com.Image.Medium)
	assert.Equal(t, src[1].Image.Small, com.Image.Small)
	assert.Empty(t, com.Genres)

	res = com.SetFromDetails(
		src,
		map[int]*genre.Genre{
			1: &genre.Genre{Name: "1"},
			2: &genre.Genre{Name: "2"},
			3: &genre.Genre{Name: "3"},
		},
	)
	assert.True(t, res)
	assert.Len(t, com.Genres, len(src[1].Genres))
}
