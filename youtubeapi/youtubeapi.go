package youtubeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Youtube ...
type Youtube struct {
	Key string
}

// New Returns a Youtube struct with the supplied key
//		key: Google API key
func New(key string) *Youtube {
	return &Youtube{
		Key: key,
	}
}

// SearchResult is the json data retrieved from the youtube api search endpoint.
type SearchResult struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []Item `json:"items"`
}

// Item stores information about videos
type Item struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	ID   struct {
		Kind    string `json:"kind"`
		VideoID string `json:"videoId"`
	} `json:"id"`
	Snippet struct {
		PublishedAt time.Time `json:"publishedAt"`
		ChannelID   string    `json:"channelId"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Thumbnails  struct {
			Default struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"default"`
			Medium struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"medium"`
			High struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"high"`
		} `json:"thumbnails"`
		ChannelTitle         string `json:"channelTitle"`
		LiveBroadcastContent string `json:"liveBroadcastContent"`
	} `json:"snippet"`
}

// Search searches youtube for videos with the supplied query.
//		query: The query to search for.
func (y *Youtube) Search(query string, maxResults int) (*SearchResult, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&key=%s&maxResults=%d", url.QueryEscape(query), y.Key, maxResults))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}

	var searchRes SearchResult
	err = json.NewDecoder(resp.Body).Decode(&searchRes)
	if err != nil {
		return nil, err
	}

	return &searchRes, nil
}

// ScrapeSearch search youtube without an api key
//		query: The query to search for.
	resp, err := http.Get("https://www.youtube.com/results?search_query=" + url.QueryEscape(query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body string
	if b, err := ioutil.ReadAll(resp.Body); err == nil {
		body = string(b)
	} else {
		return nil, err
	}

	videostart := `<a href="/watch?v=`
	videoEnd := `"`

	urls := []string{}
		startIndex := strings.Index(body, videostart)
		if startIndex < 0 {
			break
		}
		startIndex += 9

		endIndex := strings.Index(body[startIndex:], videoEnd)
		if endIndex < 0 {
			break
		}
		endIndex += startIndex

		urls = append(urls, "https://www.youtube.com"+body[startIndex:endIndex])
		body = body[endIndex:]
	}

	return urls, nil
}
