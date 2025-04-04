package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

var TrackedAccountsCommand = &discordgo.ApplicationCommand{
	Name:        "tracking",
	Description: "accounts you're currently tracking",
}

func HandleTrackedAccountsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer the response immediately
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		utils.LogError("Error deferring response", err)
		return
	}

	accounts, err := data.TrackedAccounts(i.GuildID, i.Member.User.ID)
	if err != nil {
		utils.EditResponseError(s, i, err)
		return
	}

	if len(accounts) == 0 {
		utils.EditResponseMessage(s, i, "You have no tracked accounts at the moment. Use `/track` to start tracking one!")
		return
	}

	currentCompetition := data.GetCurrentBoss(i.GuildID)
	description := ""

	if len(currentCompetition) == 0 {
		description = "\n"
	} else {
		description = fmt.Sprintf("**Event:** %s\n\n", currentCompetition)
	}

	for _, account := range accounts {
		if len(currentCompetition) > 0 {
			activity, ok := account.Activities[currentCompetition]
			if ok {
				description += fmt.Sprintf(
					"🔹 **%s**\n   └ **KC**: `%d`\n\n",
					account.Name,
					activity.CurrentAmount-activity.StartAmount,
				)
			} else {
				description += fmt.Sprintf("🔹 **%s**\n   └ *Not participating in the current event*\n", account.Name)
			}
		} else {
			description += fmt.Sprintf("🔹 **%s**\n", account.Name)
		}
	}

	embed := &discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://runetracker.org/skills/overall.gif",
		},
		Title:       "Currently Tracked Accounts",
		Description: description,
		Color:       0x00ffcc,
	}

	// Edit the deferred response with the embed
	embeds := []*discordgo.MessageEmbed{embed}
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})
	if err != nil {
		utils.LogError("Error editing response", err)
	}
}
