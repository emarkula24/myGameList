package service

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"example.com/mygamelist/repository"
)

// Defines a giantbomb service controller.
type GiantBombClient struct{}

// GetGames returns game search data from Giantbomb APi.
func (c *GiantBombClient) SearchGames(query string) (*http.Response, error) {
	apiKey := os.Getenv("GIANT_BOMB_API_KEY")
	if apiKey == "" {
		log.Println("GIANTBOMB_API_KEY is not set")
	}
	url := "https://www.giantbomb.com/api/search/?api_key=" + apiKey + "&format=json&query=" + query + "&resources=game&limit=50"
	return http.Get(url)
}

// GetGame returns game data from Giantbomb API.
func (c *GiantBombClient) SearchGame(guid string) (*http.Response, error) {
	apiKey := os.Getenv("GIANT_BOMB_API_KEY")
	url := "https://www.giantbomb.com/api/game/" + guid + "/?api_key=" + apiKey + "&format=json"
	return http.Get(url)
}

// GetGameList returns the gamedata for a users list.
func (c *GiantBombClient) SearchGameList(games []repository.Game, limit int) (*http.Response, error) {
	apiKey := os.Getenv("GIANT_BOMB_API_KEY")
	baseURL := "https://www.giantbomb.com/api/games/"
	fields := "id,guid,name,original_release_date,image"
	var ids []string
	for _, game := range games {
		idStr := strconv.Itoa(game.GameID)
		ids = append(ids, idStr)
	}

	filter := BuildUrl(ids, limit)
	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("format", "json")
	params.Set("field_list", fields)
	params.Set("filter", filter)

	url := baseURL + "?" + params.Encode()
	return http.Get(url)
}

// BuildUrl creates a query string with user gamelist id content.
func BuildUrl(ids []string, limit int) string {
	if len(ids) == 0 {
		return ""
	}
	if limit <= 0 {
		limit = len(ids)
	}
	// Return the first N IDs joined by pipe.
	chunk := ids
	if len(ids) > limit {
		chunk = ids[:limit]
	}
	filter := "id:" + strings.Join(chunk, "|")
	return filter
}
