package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/domain"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

var TrackedAccountsCommand = &discordgo.ApplicationCommand{
	Name:        "tracked-accounts",
	Description: "View all your tracked OSRS accounts",
}

func HandleTrackedAccountsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := deferResponse(s, i, "tracked-accounts"); err != nil {
		return
	}

	accounts, err := data.ParticipantRepo.GetTrackedAccounts(i.GuildID, i.Member.User.ID)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get tracked accounts")
		return
	}

	if len(accounts) == 0 {
		description := "You don't have any accounts tracked yet.\nUse `/add-account` to start tracking your OSRS accounts!"
		embed := createInfoEmbed("📋 Your Tracked Accounts", description, i)
		sendEmbedResponse(s, i, embed)
		return
	}

	hasActiveBotm, err := data.CompetitionRepo.HasRunningBotmCompetition(i.GuildID)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get competition status")
		return
	}

	var currentCompetition string
	var currentBotm *domain.BotmWithActivity
	if hasActiveBotm {
		botm, err := data.CompetitionRepo.GetBotm(i.GuildID)
		if err != nil {
			handleCommandError(s, i, err, "Failed to get competition details")
			return
		}
		if botm != nil && botm.Activity != nil {
			currentCompetition = botm.Activity.Name
			currentBotm = botm
		}
	}

	var description string
	if len(currentCompetition) == 0 {
		description = ""
	} else {
		description = fmt.Sprintf("**Event:** %s\n\n", currentCompetition)
	}

	if len(currentCompetition) > 0 {
		// Get all accounts with KC in one efficient call
		participant, err := data.ParticipantRepo.GetParticipantWithAccountKC(i.Member.User.ID, currentBotm.ID)
		if err != nil {
			utils.Error("Failed to get participant with account KC: %v", err)
			// Fallback to showing accounts without KC
			for _, account := range accounts {
				description += fmt.Sprintf("🔹 **%s**\n   └ **KC**: `Error`\n\n", account)
			}
		} else {
			// Use the KC data from the efficient query
			for _, account := range participant.Accounts {
				description += fmt.Sprintf(
					"🔹 **%s**\n   └ **KC**: `%d`\n\n",
					account.Username,
					account.KCGained,
				)
			}
		}
	} else {
		// No active competition, just show account names
		for _, account := range accounts {
			description += fmt.Sprintf("🔹 **%s**\n", account)
		}
	}

	description += "\nUse `/remove-account` to stop tracking an account."

	embed := createInfoEmbed("📋 Your Tracked Accounts", description, i)
	sendEmbedResponse(s, i, embed)
}
