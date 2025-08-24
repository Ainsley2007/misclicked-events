package commands

import (
	"fmt"
	"misclicked-events/internal/data"

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
	if hasActiveBotm {
		botm, err := data.CompetitionRepo.GetBotm(i.GuildID)
		if err != nil {
			handleCommandError(s, i, err, "Failed to get competition details")
			return
		}
		if botm != nil && botm.Activity != nil {
			currentCompetition = botm.Activity.Name
		}
	}

	var description string
	if len(currentCompetition) == 0 {
		description = ""
	} else {
		description = fmt.Sprintf("**Event:** %s\n\n", currentCompetition)
	}

	for _, account := range accounts {
		if len(currentCompetition) > 0 {
			description += fmt.Sprintf(
				"🔹 **%s**\n   └ **KC**: `%d`\n\n",
				account,
				0, // Placeholder KC value for now
			)
		} else {
			description += fmt.Sprintf("🔹 **%s**\n", account)
		}
	}

	description += "\nUse `/remove-account` to stop tracking an account."

	embed := createInfoEmbed("📋 Your Tracked Accounts", description, i)
	sendEmbedResponse(s, i, embed)
}
