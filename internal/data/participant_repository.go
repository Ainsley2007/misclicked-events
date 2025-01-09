package data

import (
	"fmt"
	"misclicked-events/internal/constants"
	"misclicked-events/internal/service"
	"misclicked-events/internal/utils"
	"sort"

	"golang.org/x/text/cases"
)

type Participant struct {
	DiscordId          string
	Points             int
	LinkedOSRSAccounts map[string]OSRSAccount
}

type OSRSAccount struct {
	Name       string
	Activities map[string]OSRSActivity
}

type OSRSActivity struct {
	Name          string
	StartAmount   int
	CurrentAmount int
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

func TrackAccount(guildID, username, discordId string) bool {
	participants, err := getParticipants(guildID)
	if err != nil {
		utils.LogError("error getting participants", err)
		return false
	}

	currentBoss := GetCurrentBoss(guildID)

	if participants == nil {
		participants = map[string]Participant{}
	}

	participant, exists := participants[discordId]

	if !exists {
		accountExists := service.CheckIfPlayerExists(username)
		if !accountExists {
			utils.LogError("err", fmt.Errorf("couldn't find a user with the name: %s", username))
			return false
		}
		participant, err = createNewParticipant(discordId, username, currentBoss)
		if err != nil {
			utils.LogError("error creating participant", err)
			return false
		}
		participants[discordId] = participant
	} else {
		err = addAccountToParticipant(&participant, username, currentBoss)
		if err != nil {
			utils.LogError("error adding account to participant", err)
			return false
		}
		participants[discordId] = participant
	}

	err = saveParticipantsData(guildID, participants)
	if err != nil {
		utils.LogError("error saving participants: ", err)
	}
	return true
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

func GetParticipantsByKCThreshold(guildID string, threshold int) ([]ParticipantKC, error) {
	participants, err := getParticipants(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch participants: %w", err)
	}

	var result []ParticipantKC

	// Iterate through participants
	for _, participant := range participants {
		totalKC := 0
		var accountBreakdown []AccountKC

		// Iterate through accounts to calculate KC
		for _, account := range participant.LinkedOSRSAccounts {
			accountKC := 0

			for _, activity := range account.Activities {
				kcDiff := activity.CurrentAmount - activity.StartAmount
				if kcDiff > threshold {
					accountKC += kcDiff
				}
			}

			// Include account only if it contributed KC
			if accountKC > 0 {
				accountBreakdown = append(accountBreakdown, AccountKC{
					AccountName: account.Name,
					TotalKC:     accountKC,
				})
				totalKC += accountKC
			}
		}

		// Add participant if they have eligible KC
		if totalKC > 0 {
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

// SortParticipantsByPointsAndKC orders participants by points descending,
// and if points are equal to 0, by total KC of their accounts descending.
func SortParticipantsByPointsAndKC(participants []Participant) []Participant {
	sort.SliceStable(participants, func(i, j int) bool {
		// Compare by points first
		if participants[i].Points != participants[j].Points {
			return participants[i].Points > participants[j].Points
		}

		// If points are equal, compare by total KC
		return getTotalKC(participants[i]) > getTotalKC(participants[j])
	})

	return participants
}

// getTotalKC calculates the total KC of all accounts for a participant
func getTotalKC(participant Participant) int {
	totalKC := 0
	for _, account := range participant.LinkedOSRSAccounts {
		for _, activity := range account.Activities {
			totalKC += activity.CurrentAmount - activity.StartAmount
		}
	}
	return totalKC
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
	for _, activityName := range constants.Activities[bossId] {
		if activity, exists := activities[activityName]; exists {
			kc += activity.Amount
		}
	}

	return max(0, kc), nil
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
