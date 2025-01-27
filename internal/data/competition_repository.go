package data

import (
	"fmt"
	"misclicked-events/internal/constants"
	"misclicked-events/internal/service"
	"misclicked-events/internal/utils"
	"sync"
)

func StartCompetition(guildID string, bossId string, competitionPassword string) error {

	saveCompetitionData(guildID, bossId, competitionPassword)

	lookupInitialKcForParticipantsAsync(guildID, bossId)

	return nil
}

func lookupInitialKcForParticipantsAsync(guildID, bossId string) {
	participants, err := getParticipants(guildID)
	if err != nil {
		fmt.Println("Error fetching participants:", err)
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex // To protect shared data (participants map)

	// Iterate through each participant
	for discordId, participant := range participants {
		// Iterate through each OSRS account linked to the participant
		for accountName, account := range participant.LinkedOSRSAccounts {
			wg.Add(1) // Increment the WaitGroup counter
			go func(discordId, accountName string, account OSRSAccount) {
				defer wg.Done() // Decrement the counter when the goroutine completes

				// Fetch hiscores for the account
				_, activities, err := service.FetchHiscore(accountName)
				if err != nil {
					fmt.Printf("Error fetching hiscore for account %s: %v\n", accountName, err)
					return
				}

				kc := 0
				for _, act := range constants.Activities[bossId].BossNames {
					activity, exists := activities[act]
					if !exists {
						fmt.Printf("No activity found for boss %s for account %s\n", bossId, accountName)
						continue
					}
					kc += max(0, activity.Amount)
				}

				// Add the initial activity to the account
				mu.Lock() // Lock the map for concurrent write
				if account.Activities == nil {
					account.Activities = make(map[string]OSRSActivity)
				}
				account.Activities[bossId] = OSRSActivity{
					Name:          bossId,
					StartAmount:   kc,
					CurrentAmount: kc,
				}
				participant.LinkedOSRSAccounts[accountName] = account
				participants[discordId] = participant
				mu.Unlock() // Unlock the map
			}(discordId, accountName, account)
		}
	}

	wg.Wait() // Wait for all goroutines to complete

	// Save the updated participants
	err = saveParticipantsData(guildID, participants)
	if err != nil {
		fmt.Println("Error saving updated participants:", err)
	}
}

func EndCompetition(guildID string, competitionPassword string) error {

	competition, err := getCompetitionData(guildID)
	if err != nil {
		return err
	}

	if len(competition.CurrentBoss) < 1 {
		return fmt.Errorf("no event is currently running")
	}

	if competition.Password != competitionPassword {
		return fmt.Errorf("incorrect event password")
	}

	err = UpdateAccountsKC(guildID)
	if err != nil {
		utils.LogError("error when updating accounts", err)
		return fmt.Errorf("error when updating accounts")
	}

	err = CalculatePointsForParticipants(guildID)
	if err != nil {
		utils.LogError("error when calculating points", err)
		return fmt.Errorf("error when calculating points")
	}

	clearCompetition(guildID)

	return nil
}
