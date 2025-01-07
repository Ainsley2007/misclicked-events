package data

import (
	"encoding/json"
	"fmt"
	"os"
)

type Competition struct {
	CurrentBoss string `json:"currentBoss"`
	Password    string `json:"password"`
}

const competitionFilePath = "./assets/%s_competition.json"

func GetCurrentBoss(guildID string) string {
	competitionData, err := getCompetitionData(guildID)
	if err != nil || competitionData == nil {
		return ""
	}

	return competitionData.CurrentBoss
}

func getCompetitionData(guildID string) (*Competition, error) {
	file, err := os.Open(fmt.Sprintf(competitionFilePath, guildID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	if stat.Size() == 0 {
		return nil, nil
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var competition Competition
	err = json.Unmarshal(data, &competition)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &competition, nil
}

func clearCompetition(guildID string) {
	saveCompetitionData(guildID, "", "")
}

func saveCompetitionData(guildID string, currentBoss string, password string) error {
	data, err := json.MarshalIndent(Competition{CurrentBoss: currentBoss, Password: password}, "", "  ") // Pretty-print
	if err != nil {
		return fmt.Errorf("failed to marshal persons: %w", err)
	}

	// Create or overwrite the file
	file, err := os.Create(fmt.Sprintf(competitionFilePath, guildID))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write JSON data to the file
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
