package data

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"

	"golang.org/x/text/cases"
)

type participantDto struct {
	DiscordId          string           `json:"id"`
	Points             int              `json:"points"`
	LinkedOSRSAccounts []oSRSAccountDto `json:"accounts"`
}

type oSRSAccountDto struct {
	Name       string        `json:"name"`
	Activities []activityDto `json:"activities"`
}

type activityDto struct {
	Name          string `json:"name"`
	StartAmount   int    `json:"startAmount"`
	CurrentAmount int    `json:"currentAmount"`
}

const (
	ErrParticipantNotFound = "participant not found"
	filePath               = "./assets/%s_participants.json"
)

func getParticipants(guildID string) (map[string]Participant, error) {
	// Open the file for reading
	file, err := os.Open(fmt.Sprintf(filePath, guildID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Check if the file is empty
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	if stat.Size() == 0 {
		// If the file is empty, return an empty slice
		return nil, nil
	}

	// Read the entire file content
	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal JSON data into a slice of Participant
	var participantsDto []participantDto
	err = json.Unmarshal(data, &participantsDto)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return mapParticipantsDto(participantsDto), nil
}

func saveParticipantsData(guildID string, participants map[string]Participant) error {
	participantsDto := mapParticipantsToDto(participants)
	data, err := json.MarshalIndent(participantsDto, "", "  ") // Pretty-print
	if err != nil {
		return fmt.Errorf("failed to marshal persons: %w", err)
	}

	// Create or overwrite the file
	file, err := os.Create(fmt.Sprintf(filePath, guildID))
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

func mapParticipantsToDto(participants map[string]Participant) []participantDto {
	participantsDto := make([]participantDto, 0, len(participants))

	// Convert map values to a slice
	participantSlice := slices.Collect(maps.Values(participants))

	for _, p := range participantSlice {
		linkedAccountsDto := make([]oSRSAccountDto, 0, len(p.LinkedOSRSAccounts))

		accountSlice := slices.Collect(maps.Values(p.LinkedOSRSAccounts))
		for _, a := range accountSlice {
			activitiesDto := make([]activityDto, 0, len(a.Activities))

			activitySlice := slices.Collect(maps.Values(a.Activities))
			for _, ac := range activitySlice {
				activitiesDto = append(activitiesDto, activityDto(ac))
			}

			linkedAccountsDto = append(linkedAccountsDto, oSRSAccountDto{
				Name:       a.Name,
				Activities: activitiesDto,
			})
		}

		participantsDto = append(participantsDto, participantDto{
			DiscordId:          p.DiscordId,
			Points:             p.Points,
			LinkedOSRSAccounts: linkedAccountsDto,
		})
	}

	return participantsDto
}

func mapParticipantsDto(participantsDto []participantDto) map[string]Participant {
	participants := make(map[string]Participant) // Initialize participants map
	for _, p := range participantsDto {
		accounts := make(map[string]OSRSAccount) // Initialize accounts map
		for _, a := range p.LinkedOSRSAccounts {
			activities := make(map[string]OSRSActivity) // Initialize activities map
			for _, ac := range a.Activities {
				activities[ac.Name] = OSRSActivity(ac) // Direct conversion
			}

			accounts[cases.Fold().String(a.Name)] = OSRSAccount{
				Name:       a.Name,
				Activities: activities,
			}
		}

		participants[p.DiscordId] = Participant{
			DiscordId:          p.DiscordId,
			Points:             p.Points,
			LinkedOSRSAccounts: accounts,
		}
	}
	return participants
}
