package api

import (
	"encoding/json"
	"fmt"
	"misclicked-events/internal/utils"
	"net/http"
	"net/url"
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
	utils.Debug("Creating new HiscoreDataSource with base URL: https://secure.runescape.com/m=hiscore_oldschool")
	return &hiscoreDS{
		client:  &http.Client{},
		baseURL: "https://secure.runescape.com/m=hiscore_oldschool",
	}
}

func (ds *hiscoreDS) CheckIfPlayerExists(username string) (bool, error) {
	utils.Debug("Checking if player %s exists in OSRS", username)

	encodedUsername := url.QueryEscape(username)
	url := fmt.Sprintf("%s/index_lite.ws?player=%s", ds.baseURL, encodedUsername)
	utils.Debug("Making HTTP GET request to: %s", url)

	resp, err := ds.client.Get(url)
	if err != nil {
		utils.Error("Failed to make HTTP request to check if player %s exists: %v", username, err)
		return false, fmt.Errorf("failed to check player")
	}
	defer resp.Body.Close()

	utils.Debug("Received response for player %s: status code %d", username, resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		utils.Debug("Player %s exists in OSRS", username)
		return true, nil
	} else {
		utils.Debug("Player %s does not exist in OSRS (status code: %d)", username, resp.StatusCode)
		return false, nil
	}
}

func (ds *hiscoreDS) FetchHiscore(username string) (*HiscoreDataModel, error) {
	utils.Debug("Fetching hiscore data for player %s", username)

	encodedUsername := url.QueryEscape(username)
	url := fmt.Sprintf("%s/index_lite.json?player=%s", ds.baseURL, encodedUsername)
	utils.Debug("Making HTTP GET request to: %s", url)

	resp, err := ds.client.Get(url)
	if err != nil {
		utils.Error("Failed to make HTTP request to fetch hiscore data for player %s: %v", username, err)
		return nil, fmt.Errorf("failed to fetch hiscore")
	}
	defer resp.Body.Close()

	utils.Debug("Received response for hiscore data for player %s: status code %d", username, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		utils.Error("Received non-200 response for hiscore data for player %s: status code %d", username, resp.StatusCode)
		return nil, fmt.Errorf("failed to fetch hiscore")
	}

	var data HiscoreDataModel
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		utils.Error("Failed to decode hiscore JSON for player %s: %v", username, err)
		return nil, fmt.Errorf("failed to fetch hiscore")
	}

	utils.Debug("Successfully decoded hiscore data for player %s: %d skills, %d activities", username, len(data.Skills), len(data.Activities))
	return &data, nil
}
