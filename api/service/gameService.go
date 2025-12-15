package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"example.com/mygamelist/repository"
)

// Defines a giantbomb service controller.
type GiantBombClient struct {
	httpClient  *http.Client
	token       string
	tokenExpiry time.Time
	mu          sync.Mutex
}
type OAuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func NewGiantBombClient() *GiantBombClient {
	return &GiantBombClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *GiantBombClient) getAPIToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}

	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")
	url := "https://id.twitch.tv/oauth2/token?client_id=" + client_id + "&client_secret=" + client_secret + "&grant_type=client_credentials"
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		nil,
	)
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api auth failed: %s", resp.Status)
	}

	var tokenResp OAuthToken
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	c.token = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)

	return c.token, nil

}
func (c *GiantBombClient) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	token, err := c.getAPIToken(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", os.Getenv("CLIENT_ID"))

	return c.httpClient.Do(req)
}

// GetGames returns game search data from Giantbomb APi.
func (c *GiantBombClient) SearchGames(ctx context.Context, query string) (*http.Response, error) {

	u := "https://api.igdb.com/v4/games"
	body := strings.NewReader(
		fmt.Sprintf(`search "%s"; fields name, cover.url, platforms.abbreviation; limit 100; 
		where cover != null & platforms.abbreviation != null & version_parent = null; `, query),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return c.doRequest(ctx, req)
}

// GetGame returns game data from Giantbomb API.
func (c *GiantBombClient) SearchGame(ctx context.Context, guid string) (*http.Response, error) {
	u := "https://api.igdb.com/v4/games"
	body := strings.NewReader(
		fmt.Sprintf(`fields cover.url, summary, name, release_dates.human, platforms.abbreviation, 
		similar_games, franchise.*, genres.name, rating, rating_count; 
		where id = %s & platforms.platform_logo.url != null;`, guid),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return c.doRequest(ctx, req)
}

// GetGameList returns the gamedata for a users list.
func (c *GiantBombClient) SearchGameList(ctx context.Context, games []repository.Game, limit int) (*http.Response, error) {
	var ids []string
	for _, game := range games {
		idStr := strconv.Itoa(game.GameID)
		ids = append(ids, idStr)
	}
	ids_string := strings.Join(ids, ",")
	u := "https://api.igdb.com/v4/games"
	body := strings.NewReader(
		fmt.Sprintf(`fields cover.url, name, release_dates.human; where id = (%s) & cover.url != null;`, ids_string),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return c.doRequest(ctx, req)
}
