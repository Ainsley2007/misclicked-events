package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HiscoreDataSource defines the interface for hiscore data operations
type HiscoreDataSource interface {
	CheckIfPlayerExists(username string) (bool, error)
	FetchHiscore(username string) (*HiscoreData, error)
}

// HiscoreData represents the complete hiscore data structure
type HiscoreData struct {
	Skills     []Skill    `json:"skills"`
	Activities []Activity `json:"activities"`
}

// Skill represents a skill in the hiscore data
type Skill struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Level int    `json:"level"`
	XP    int    `json:"xp"`
}

// Activity represents an activity/boss in the hiscore data
type Activity struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Score int    `json:"score"`
}

// hiscoreDS implements HiscoreDataSource
type hiscoreDS struct {
	client  *http.Client
	baseURL string
}

// NewHiscoreDataSource creates a new instance of HiscoreDataSource
func NewHiscoreDataSource() HiscoreDataSource {
	return &hiscoreDS{
		client:  &http.Client{},
		baseURL: "https://secure.runescape.com/m=hiscore_oldschool",
	}
}

// CheckIfPlayerExists checks if a player exists by making a request to the hiscore API
func (ds *hiscoreDS) CheckIfPlayerExists(username string) (bool, error) {
	url := fmt.Sprintf("%s/index_lite.ws?player=%s", ds.baseURL, username)

	resp, err := ds.client.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to check if player exists: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// FetchHiscore fetches the complete hiscore data for a player
func (ds *hiscoreDS) FetchHiscore(username string) (*HiscoreData, error) {
	url := fmt.Sprintf("%s/index_lite.json?player=%s", ds.baseURL, username)

	resp, err := ds.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch hiscore data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var data HiscoreData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode hiscore JSON: %w", err)
	}

	return &data, nil
}
