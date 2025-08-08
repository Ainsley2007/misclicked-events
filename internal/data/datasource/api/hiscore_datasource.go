package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HiscoreDataSource interface {
	CheckIfPlayerExists(username string) (bool, error)
	FetchHiscore(username string) (*HiscoreDataModel, error)
}

type hiscoreDS struct {
	client  *http.Client
	baseURL string
}

func NewHiscoreDataSource() HiscoreDataSource {
	return &hiscoreDS{
		client:  &http.Client{},
		baseURL: "https://secure.runescape.com/m=hiscore_oldschool",
	}
}

func (ds *hiscoreDS) CheckIfPlayerExists(username string) (bool, error) {
	url := fmt.Sprintf("%s/index_lite.ws?player=%s", ds.baseURL, username)

	resp, err := ds.client.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to check if player exists: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

func (ds *hiscoreDS) FetchHiscore(username string) (*HiscoreDataModel, error) {
	url := fmt.Sprintf("%s/index_lite.json?player=%s", ds.baseURL, username)

	resp, err := ds.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch hiscore data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var data HiscoreDataModel
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode hiscore JSON: %w", err)
	}

	return &data, nil
}
