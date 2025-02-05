package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Skill represents a single skill object.
type Skill struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Level int    `json:"level"`
	XP    int    `json:"xp"`
}

// Activity represents a single activity object.
type Activity struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Score int    `json:"score"`
}

func CheckIfPlayerExists(username string) bool {
	url := fmt.Sprintf("https://secure.runescape.com/m=hiscore_oldschool/index_lite.ws?player=%s", username)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}

	return false
}

func FetchHiscore(username string) ([]Skill, []Activity, error) {
	url := fmt.Sprintf("https://secure.runescape.com/m=hiscore_oldschool/index_lite.json?player=%s", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch hiscore data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Temporary struct matching the JSON structure.
	var data struct {
		Skills     []Skill    `json:"skills"`
		Activities []Activity `json:"activities"`
	}

	// Decode the JSON response.
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil, fmt.Errorf("failed to decode hiscore JSON: %w", err)
	}

	return data.Skills, data.Activities, nil
}

// FindSkill searches for a skill by its name in the provided slice.
// It returns a pointer to the found Skill and true if found; otherwise, nil and false.
func FindSkill(skills []Skill, name string) (*Skill, bool) {
	for i, skill := range skills {
		if skill.Name == name {
			return &skills[i], true
		}
	}
	return nil, false
}

// FindActivity searches for an activity by its name in the provided slice.
// It returns a pointer to the found Activity and true if found; otherwise, nil and false.
func FindActivity(activities []Activity, name string) (*Activity, bool) {
	for i, activity := range activities {
		if activity.Name == name {
			return &activities[i], true
		}
	}
	return nil, false
}
