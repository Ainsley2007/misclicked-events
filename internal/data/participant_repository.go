package data

import (
	"fmt"
	"maps"
	"misclicked-events/internal/constants"
	"misclicked-events/internal/service"
	"misclicked-events/internal/utils"
	"slices"
	"sort"

	"golang.org/x/text/cases"
)

type Participant struct {
	DiscordId          string
	Points             int
	LinkedOSRSAccounts map[string]OSRSAccount
}

func (p Participant) TotalKCForActivity(activityName string) (int, []AccountKC) {
	totalKC := 0
	var accountBreakdown []AccountKC

	for _, account := range p.LinkedOSRSAccounts {
		accountKC := account.KCForActivity(activityName)
		if accountKC > 0 { // Include accounts contributing more than 0 KC
			accountBreakdown = append(accountBreakdown, AccountKC{
				AccountName: account.Name,
				TotalKC:     accountKC,
			})
			totalKC += accountKC
		}
	}

	// Sort the breakdown by KC in descending order
	sort.Slice(accountBreakdown, func(i, j int) bool {
		return accountBreakdown[i].TotalKC > accountBreakdown[j].TotalKC
	})

	return totalKC, accountBreakdown
}

type OSRSAccount struct {
	Name       string
	Activities map[string]OSRSActivity
}

func (acc OSRSAccount) KCForActivity(activityName string) int {
	activity, exists := acc.Activities[activityName]
	if !exists {
		return 0
	}
	return activity.KC()
}

type OSRSActivity struct {
	Name          string
	StartAmount   int
	CurrentAmount int
}

func (a OSRSActivity) KC() int {
	return a.CurrentAmount - a.StartAmount
}

type ParticipantKC struct {
	DiscordId  string
	TotalKC    int
	AccountKCs []AccountKC
}

type AccountKC struct {
	AccountName string
	TotalKC     int
}

func TrackAccount(guildID, username, discordId string) error {
	// Retrieve participants for the guild
	participants, err := getParticipants(guildID)
	if err != nil {
		utils.LogError("Failed to retrieve participants", err)
		return fmt.Errorf("failed to retrieve participants: %w", err)
	}

	// Ensure the participants map is initialized
	if participants == nil {
		participants = make(map[string]Participant)
	}

	// Get the current competition boss (if any)
	currentBoss := GetCurrentBoss(guildID)

	// Check if the participant already exists
	participant, exists := participants[discordId]

	if !exists {
		// Validate the username
		if !service.CheckIfPlayerExists(username) {
			err := fmt.Errorf("could not find an OSRS account with the username: %s", username)
			utils.LogError("Invalid username", err)
			return err
		}

		// Create a new participant
		participant, err = createNewParticipant(discordId, username, currentBoss)
		if err != nil {
			utils.LogError("Failed to create new participant", err)
			return fmt.Errorf("failed to create new participant: %w", err)
		}

		participants[discordId] = participant
	} else {
		// Validate the username
		if !service.CheckIfPlayerExists(username) {
			err := fmt.Errorf("could not find an OSRS account with the username: %s", username)
			utils.LogError("Invalid username", err)
			return err
		}

		// Add the account to the existing participant
		err = addAccountToParticipant(&participant, username, currentBoss)
		if err != nil {
			utils.LogError("Failed to add account to participant", err)
			return fmt.Errorf("failed to add account to participant: %w", err)
		}

		participants[discordId] = participant
	}

	// Save the updated participants map
	err = saveParticipantsData(guildID, participants)
	if err != nil {
		utils.LogError("Failed to save participants data", err)
		return fmt.Errorf("failed to save participants data: %w", err)
	}

	return nil
}

func UpdateAccountsKC(guildID string) error {
	// Fetch all participants
	participants, err := getParticipants(guildID)
	if err != nil {
		return fmt.Errorf("failed to fetch participants: %w", err)
	}

	// Get the currently ongoing boss
	currentBoss := GetCurrentBoss(guildID)
	if currentBoss == "" {
		return fmt.Errorf("no ongoing boss competition")
	}

	// Iterate through each participant
	for discordId, participant := range participants {
		updated := false

		// Iterate through each linked OSRS account
		for username, account := range participant.LinkedOSRSAccounts {
			// Fetch the current KC for the account and boss
			kc, err := fetchKc(username, currentBoss)
			if err != nil {
				fmt.Printf("Error fetching KC for account %s: %v\n", username, err)
				continue
			}

			// Update the activity if it exists
			if account.Activities == nil {
				account.Activities = make(map[string]OSRSActivity)
			}

			if activity, exists := account.Activities[currentBoss]; exists {
				activity.CurrentAmount = kc
				account.Activities[currentBoss] = activity
			} else {
				// Add a new activity if not already tracked
				account.Activities[currentBoss] = OSRSActivity{
					Name:          currentBoss,
					StartAmount:   kc,
					CurrentAmount: kc,
				}
			}

			participant.LinkedOSRSAccounts[username] = account
			updated = true
		}

		// Update the participant only if something was updated
		if updated {
			participants[discordId] = participant
		}
	}

	// Save updated participants
	err = saveParticipantsData(guildID, participants)
	if err != nil {
		return fmt.Errorf("failed to save updated participants: %w", err)
	}

	return nil
}

func GetParticipantsByActivityKCThreshold(guildID string) ([]ParticipantKC, error) {
	// Fetch participants for the given guild
	participants, err := getParticipants(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch participants: %w", err)
	}

	activityName := GetCurrentBoss(guildID)
	if activityName == "" {
		return nil, fmt.Errorf("no event found")
	}

	var result []ParticipantKC

	// Iterate through participants
	for _, participant := range participants {
		// Use the Participant method to calculate total KC and breakdown
		totalKC, accountBreakdown := participant.TotalKCForActivity(activityName)

		// Add participant to the result if their total KC exceeds the threshold
		if totalKC >= constants.Activities[activityName].Threshold {
			result = append(result, ParticipantKC{
				DiscordId:  participant.DiscordId,
				TotalKC:    totalKC,
				AccountKCs: accountBreakdown,
			})
		}
	}

	// Sort the result by TotalKC in descending order
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalKC > result[j].TotalKC
	})

	return result, nil
}

// createNewParticipant initializes a new participant with the given username
// and optional boss KC if a boss is active.
func createNewParticipant(discordId, username, currentBoss string) (Participant, error) {
	activities := map[string]OSRSActivity{}

	// Fetch initial KC if a boss is active
	if currentBoss != "" {
		kc, err := fetchKc(username, currentBoss)
		if err != nil {
			return Participant{}, err
		}
		activities[currentBoss] = OSRSActivity{
			Name:          currentBoss,
			StartAmount:   kc,
			CurrentAmount: kc,
		}
	}

	return Participant{
		DiscordId: discordId,
		Points:    0,
		LinkedOSRSAccounts: map[string]OSRSAccount{
			username: {
				Name:       username,
				Activities: activities,
			},
		},
	}, nil
}

// addAccountToParticipant adds a new account to an existing participant.
// If a boss is active, it initializes KC for the current boss.
func addAccountToParticipant(participant *Participant, username, currentBoss string) error {
	usernameKey := cases.Fold().String(username)
	if _, exists := participant.LinkedOSRSAccounts[usernameKey]; exists {
		return fmt.Errorf("account is already being tracked")
	}

	activities := map[string]OSRSActivity{}

	// Fetch initial KC if a boss is active
	if currentBoss != "" {
		kc, err := fetchKc(username, currentBoss)
		if err != nil {
			return err
		}
		activities[currentBoss] = OSRSActivity{
			Name:          currentBoss,
			StartAmount:   kc,
			CurrentAmount: kc,
		}
	}

	participant.LinkedOSRSAccounts[usernameKey] = OSRSAccount{
		Name:       username,
		Activities: activities,
	}

	return nil
}

// fetchKc calculates the total KC for the given username and boss.
func fetchKc(username, bossId string) (int, error) {
	_, activities, err := service.FetchHiscore(username)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch hiscores for %s: %w", username, err)
	}

	kc := 0
	for _, activityName := range constants.Activities[bossId].BossNames {
		if activity, exists := activities[activityName]; exists {
			kc += max(0, activity.Amount)
		}
	}

	return max(0, kc), nil
}

func GetParticipantsInOrder(guildID string) ([]Participant, error) {
	participants, err := getParticipants(guildID)
	if err != nil {
		return []Participant{}, err
	}

	parts := slices.Collect(maps.Values(participants))

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Points > parts[j].Points
	})

	return parts, nil
}

func UntrackAccount(guildID, username, discordId string) error {
	participants, err := getParticipants(guildID)
	if err != nil {
		return err
	}

	participant, ok := participants[discordId]

	if !ok {
		return fmt.Errorf("we are currently not tracking any accounts for you")
	}

	usernameKey := cases.Fold().String(username)
	_, ok = participant.LinkedOSRSAccounts[usernameKey]

	if !ok {
		return fmt.Errorf("no account found by this name")
	}

	delete(participant.LinkedOSRSAccounts, usernameKey)

	if len(participant.LinkedOSRSAccounts) == 0 {
		delete(participants, discordId)
	} else {
		participants[discordId] = participant
	}

	return saveParticipantsData(guildID, participants)
}

func TrackedAccounts(guildId, discordId string) ([]OSRSAccount, error) {
	participants, err := getParticipants(guildId)
	if err != nil {
		return nil, err
	}

	participant, exists := participants[discordId]
	if !exists {
		return nil, fmt.Errorf("we are currently not tracking any accounts for you")
	}

	accounts := make([]OSRSAccount, 0, len(participant.LinkedOSRSAccounts))
	for _, account := range participant.LinkedOSRSAccounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// CalculatePointsForParticipants calculates and assigns points to participants based on their TotalKC.
func CalculatePointsForParticipants(guildID string) error {
	// Get participants above the threshold, sorted by TotalKC (descending)
	participantsAboveThreshold, err := GetParticipantsByActivityKCThreshold(guildID)
	if err != nil {
		return fmt.Errorf("failed to get participants above the threshold: %w", err)
	}

	// Retrieve the full participants map
	participants, err := getParticipants(guildID)
	if err != nil {
		return fmt.Errorf("failed to retrieve participants: %w", err)
	}

	// Define the point system
	pointSystem := []int{12, 9, 7, 5, 4, 3, 3, 2, 2, 2} // Points for ranks 1 through 10

	// Track current rank and previous KC
	currentRank := 1
	previousKC := -1
	pointsToAward := 0

	// Iterate through the participants and calculate points
	for i, participantKC := range participantsAboveThreshold {
		// If the current participant's TotalKC differs from the previous, update rank and points
		if participantKC.TotalKC != previousKC {
			currentRank = i + 1

			// Determine points to award based on rank
			if currentRank <= len(pointSystem) {
				pointsToAward = pointSystem[currentRank-1]
			} else {
				pointsToAward = 1 // Default point for ranks beyond the defined system
			}
		}

		// Update the points for the participant in the map
		if participant, exists := participants[participantKC.DiscordId]; exists {
			participant.Points += pointsToAward
			participants[participantKC.DiscordId] = participant
		} else {
			return fmt.Errorf("participant with Discord ID %s not found in map", participantKC.DiscordId)
		}

		// Update the previous KC to the current one
		previousKC = participantKC.TotalKC
	}

	// Save the updated participants map
	err = saveParticipantsData(guildID, participants)
	if err != nil {
		return fmt.Errorf("failed to save updated participants: %w", err)
	}

	return nil
}
