package data

import (
	"fmt"
	"sync"
)

func StartCompetition(guildID string, bossId string, competitionPassword string) error {

	saveCompetitionData(guildID, bossId, competitionPassword)

	lookupInitialKcForParticipantsAsync(guildID, bossId)

	return nil
}

/*
	func lookupInitialKcForParticipants(bossId string) {
		// Fetch all participants
		participants, err := getParticipants()
		if err != nil {
			fmt.Println("Error fetching participants:", err)
			return
		}

		// Iterate through each participant
		for discordId, participant := range participants {
			// Iterate through each OSRS account linked to the participant
			for accountName, account := range participant.LinkedOSRSAccounts {
				// Fetch hiscores for the account
				_, activities, err := fetchHiscore(accountName)
				if err != nil {
					fmt.Printf("Error fetching hiscore for account %s: %v\n", accountName, err)
					continue
				}

				kc := 0
				for _, act := range Activities[bossId] {
					activity, exists := activities[act]
					if !exists {
						fmt.Printf("No activity found for boss %s for account %s\n", bossId, accountName)
						continue
					}
					kc += activity.Amount
				}

				// Add the initial activity to the account
				if account.Activities == nil {
					account.Activities = make(map[string]OSRSActivity)
				}
				account.Activities[bossId] = OSRSActivity{
					Name:          bossId,
					StartAmount:   kc,
					CurrentAmount: kc,
				}

				// Update the participant's account in the map
				participant.LinkedOSRSAccounts[accountName] = account
			}

			// Update the participants map
			participants[discordId] = participant
		}

		// Save the updated participants
		err = saveParticipantsData(participants)
		if err != nil {
			fmt.Println("Error saving updated participants:", err)
		}
	}
*/
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
				_, activities, err := fetchHiscore(accountName)
				if err != nil {
					fmt.Printf("Error fetching hiscore for account %s: %v\n", accountName, err)
					return
				}

				kc := 0
				for _, act := range Activities[bossId] {
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

	if competition.Password == competitionPassword {
		clearCompetition(guildID)
	} else {
		return fmt.Errorf("incorrect event password")
	}

	return nil
}
