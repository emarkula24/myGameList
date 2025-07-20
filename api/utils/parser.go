package utils

import (
	"encoding/json"
	"strings"

	"example.com/mygamelist/repository"
)

func ParseSearchQuery(query string) string {
	querysplice := strings.ReplaceAll(query, " ", "-")
	return querysplice
}

type ApiImage struct {
	IconURL     string `json:"icon_url"`
	MediumURL   string `json:"medium_url"`
	ScreenURL   string `json:"screen_url"`
	ScreenLarge string `json:"screen_large_url"`
	SmallURL    string `json:"small_url"`
	SuperURL    string `json:"super_url"`
	ThumbURL    string `json:"thumb_url"`
	TinyURL     string `json:"tiny_url"`
	OriginalURL string `json:"original_url"`
	ImageTags   string `json:"image_tags"`
}

type ApiResult struct {
	ID                  int      `json:"id"`
	Image               ApiImage `json:"image"`
	Name                string   `json:"name"`
	OriginalReleaseDate string   `json:"original_release_date"`
	Status              string   `json:"status,omitempty"` // this value is injected into the ApiResponse
}

type ApiResponse struct {
	Error                string      `json:"error"`
	Limit                int         `json:"limit"`
	Offset               int         `json:"offset"`
	NumberOfPageResults  int         `json:"number_of_page_results"`
	NumberOfTotalResults int         `json:"number_of_total_results"`
	StatusCode           int         `json:"status_code"`
	Results              []ApiResult `json:"results"`
	Version              string      `json:"version"`
}

// Builds a struct that is used in response data, which contains status value from database.
func CombineGameListJSON(gameListDb []repository.Game, bodyBytes []byte) (*ApiResponse, error) {
	var apiResp ApiResponse
	err := json.Unmarshal(bodyBytes, &apiResp)
	if err != nil {
		return nil, err
	}
	statusMap := make(map[int]string)
	for _, g := range gameListDb {
		statusMap[g.GameID] = g.Status
	}

	for i := range apiResp.Results {
		if status, ok := statusMap[apiResp.Results[i].ID]; ok {
			apiResp.Results[i].Status = status
		}
	}
	return &apiResp, nil
}
