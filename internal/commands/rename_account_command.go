package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var RenameAccountCommand = &discordgo.ApplicationCommand{
	Name:        "rename-account",
	Description: "Rename one of your tracked OSRS accounts",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "old_username",
			Description:  "The current username of the account",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "new_username",
			Description: "The new username to change to",
			Required:    true,
		},
	},
}

func HandleRenameAccountCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := deferResponse(s, i, "rename-account"); err != nil {
		return
	}

	if err := validateRequiredOptions(i, 2, "rename-account"); err != nil {
		handleCommandError(s, i, err, "Rename account command validation failed")
		return
	}

	oldUsername, err := getStringOption(i, 0)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get old username option")
		return
	}

	newUsername, err := getStringOption(i, 1)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get new username option")
		return
	}

	err = data.RenameAccountUseCase.Execute(i.GuildID, i.Member.User.ID, oldUsername, newUsername)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Player Not Found",
				Description: "The new player does not exist in OSRS.",
				Color:       0xff0000,
			}
			utils.EditResponseEmbed(s, i, embed)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			embed := &discordgo.MessageEmbed{
				Title:       "❌ Account Not Found",
				Description: "The old account is not being tracked.",
				Color:       0xff0000,
			}
			utils.EditResponseEmbed(s, i, embed)
			return
		}
		handleCommandError(s, i, err, "Failed to rename account")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "✅ Account Renamed Successfully",
		Description: fmt.Sprintf("Successfully renamed account **%s** → **%s**\n\nYour account has been updated in all competitions.", oldUsername, newUsername),
		Color:       0x00ff00,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", i.Member.User.Username),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	utils.EditResponseEmbed(s, i, embed)
}
