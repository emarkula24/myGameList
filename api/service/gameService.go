package service

import (
	"log"
	"net/http"
	"os"
)

type GiantBombClient struct{}

func (c *GiantBombClient) SearchGames(query string) (*http.Response, error) {
	apiKey := os.Getenv("GIANT_BOMB_API_KEY")
	if apiKey == "" {
		log.Println("GIANTBOMB_API_KEY is not set")
	}
	url := "https://www.giantbomb.com/api/search/?api_key=" + apiKey + "&format=json&query=" + query + "&resources=game&limit=50"
	return http.Get(url)
}

func (c *GiantBombClient) SearchGame(guid string) (*http.Response, error) {
	apiKey := os.Getenv("GIANT_BOMB_API_KEY")
	url := "https://www.giantbomb.com/api/game/" + guid + "/?api_key=" + apiKey + "&format=json"
	return http.Get(url)
}
