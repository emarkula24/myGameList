package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"example.com/mygamelist/repository"
)

func ParseSearchQuery(query string) string {
	querysplice := strings.ReplaceAll(query, " ", "-")
	return querysplice
}

type ApiResult struct {
	ID                  int    `json:"id"`
	Cover               string `json:"cover"`
	Name                string `json:"name"`
	OriginalReleaseDate string `json:"original_release_date"`
	Status              int    `json:"status,omitempty"` // this value is injected into the ApiResponse
}

type igdbSearchResult struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Cover struct {
		URL string `json:"url"`
	} `json:"cover"`
	ReleaseDates []struct {
		Human string `json:"human"`
	} `json:"release_dates"`
}

// Builds a struct that is used in response data, which contains status value from database.
func CombineGameListJSON(gameListDb []repository.Game, bodyBytes []byte) ([]ApiResult, error) {

	var igdbResults []igdbSearchResult
	if err := json.Unmarshal(bodyBytes, &igdbResults); err != nil {
		return nil, fmt.Errorf("igdb unmarshal failed: %w", err)
	}

	statusMap := make(map[int]int)
	for _, g := range gameListDb {
		statusMap[g.GameID] = g.Status
	}

	results := make([]ApiResult, 0, len(igdbResults))
	for _, g := range igdbResults {
		api := igdbToApiResult(g)

		// Inject status if present
		if status, ok := statusMap[g.ID]; ok {
			api.Status = status
		}

		results = append(results, api)
	}

	return results, nil
}

func igdbToApiResult(src igdbSearchResult) ApiResult {
	var release string
	if len(src.ReleaseDates) > 0 {
		release = src.ReleaseDates[len(src.ReleaseDates)-1].Human
	}

	return ApiResult{
		ID:                  src.ID,
		Name:                src.Name,
		Cover:               src.Cover.URL,
		OriginalReleaseDate: release,
	}
}
