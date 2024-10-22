package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

type TMDBClient struct {
	apiKey string
	client *http.Client
}

type TMDBMovie struct {
	Title         string   `json:"title"`
	OriginalTitle string   `json:"original_title"`
	PosterPath    string   `json:"poster_path"`
	Genres        []genre  `json:"genres"`
	OriginCountry []string `json:"origin_country"`
	ReleaseDate   string   `json:"release_date"`
}

type genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewTMDBClient(apiKey string) *TMDBClient {
	return &TMDBClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (c *TMDBClient) GetMovie(id int) (TMDBMovie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?api_key=%s&language=pt-BR", id, c.apiKey)
	resp, err := c.client.Get(url)
	if err != nil {
		slog.Error("Failed to get movie from TMDB", "error", err)
		return TMDBMovie{}, err
	}
	defer resp.Body.Close()
	log.Println("TMDB response status:", resp.Status)

	decoder := json.NewDecoder(resp.Body)
	var movie TMDBMovie
	if err := decoder.Decode(&movie); err != nil {
		slog.Error("Failed to decode movie from TMDB", "error", err)
		return TMDBMovie{}, err
	}

	return movie, nil
}
