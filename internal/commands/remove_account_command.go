package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"
	"strings"
	"time"

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
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Account Not Found",
				Description: "This account is not being tracked.",
				Color:       0xff0000,
			}
			utils.EditResponseEmbed(s, i, embed)
			return
		}
		handleCommandError(s, i, err, "Failed to remove account")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "✅ Account Removed Successfully",
		Description: fmt.Sprintf("Successfully removed account **%s**\n\nYour account is no longer being tracked for competitions.", username),
		Color:       0x00ff00,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", i.Member.User.Username),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	utils.EditResponseEmbed(s, i, embed)
}
