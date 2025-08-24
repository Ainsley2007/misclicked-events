package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var RemoveAccountCommand = &discordgo.ApplicationCommand{
	Name:        "remove-account",
	Description: "Stop tracking an OSRS account from your profile",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "username",
			Description:  "The OSRS account username to stop tracking",
			Required:     true,
			Autocomplete: true,
		},
	},
}

func HandleRemoveAccountCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := deferResponse(s, i, "remove-account"); err != nil {
		return
	}

	if err := validateRequiredOptions(i, 1, "remove-account"); err != nil {
		handleCommandError(s, i, err, "Remove account command validation failed")
		return
	}

	username, err := getStringOption(i, 0)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get username option")
		return
	}

	err = data.ParticipantRepo.RemoveAccount(i.GuildID, i.Member.User.ID, username)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			handleAccountNotFoundError(s, i, username)
			return
		}
		handleCommandError(s, i, err, "Failed to remove account")
		return
	}

	description := fmt.Sprintf("Successfully removed account **%s**\n\nYour account is no longer being tracked for competitions.", username)
	embed := createSuccessEmbed("✅ Account Removed Successfully", description, i)
	sendEmbedResponse(s, i, embed)
}
