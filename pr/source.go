package pr

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Source struct {
	c http.Client
}

func NewSource(c http.Client) *Source {
	return &Source{c: c}
}

func (s Source) Find(ctx context.Context, target string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://apipodcasts.polskieradio.pl/api/podcasts", nil)
	if err != nil {
		return nil, err
	}
	res, err := s.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing body: %v", err)
		}
	}(res.Body)

	decoder := json.NewDecoder(res.Body)
	var podcasts Podcasts
	err = decoder.Decode(&podcasts)
	if err != nil {
		return nil, err
	}
	targets := make(map[string]string)
	for _, p := range podcasts.Items {
		targets[p.Title] = p.PodcastUrl
		targets[p.Description] = p.PodcastUrl
		targets[p.ItunesKeywords] = p.PodcastUrl
		targets[p.ItunesCategory] = p.PodcastUrl
	}
	var keys []string
	for k := range targets {
		keys = append(keys, k)
	}
	matches := fuzzy.Find(target, keys)
	var results []string
	for _, m := range matches {
		results = append(results, targets[m])
	}
	return results, nil
}

type Podcasts struct {
	Items []struct {
		Id             int       `json:"id"`
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		RadioStation   string    `json:"radioStation"`
		Audition       string    `json:"audition"`
		Language       string    `json:"language"`
		PodcastUrl     string    `json:"podcastUrl"`
		ItunesKeywords string    `json:"itunesKeywords"`
		ItunesCategory string    `json:"itunesCategory"`
		PrCategory     *string   `json:"prCategory"`
		CreationDate   time.Time `json:"creationDate"`
		Guid           string    `json:"guid"`
		ViewCount      int       `json:"viewCount"`
		ItemCount      int       `json:"itemCount"`
		LastItemUpdate time.Time `json:"lastItemUpdate"`
		Image          struct {
			Main        string `json:"main"`
			Landscape   string `json:"landscape"`
			Thumbnail   string `json:"thumbnail"`
			Recommended string `json:"recommended"`
		} `json:"image"`
		StreamingServices []struct {
			Name   string  `json:"name"`
			Link   *string `json:"link"`
			Status *string `json:"status"`
		} `json:"streamingServices"`
		Announcer    *string `json:"announcer"`
		AnnouncerImg string  `json:"announcerImg"`
		Email        *string `json:"email"`
		Socials      *struct {
			Img         string `json:"img"`
			Site        string `json:"site"`
			SocialMedia struct {
				Fb      string  `json:"fb"`
				Twitter *string `json:"twitter"`
				Yt      *string `json:"yt"`
			} `json:"socialMedia"`
		} `json:"socials"`
	} `json:"items"`
	Count    int `json:"count"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
