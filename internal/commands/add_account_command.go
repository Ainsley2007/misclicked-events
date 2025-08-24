package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var AddAccountCommand = &discordgo.ApplicationCommand{
	Name:        "add-account",
	Description: "Link an OSRS account to your profile to track its progress — only add accounts you own.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "username",
			Description: "The OSRS username to start tracking",
			Required:    true,
		},
	},
}

func HandleAddAccountCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := deferResponse(s, i, "add-account"); err != nil {
		return
	}

	if err := validateRequiredOptions(i, 1, "add-account"); err != nil {
		handleCommandError(s, i, err, "Add account command validation failed")
		return
	}

	username, err := getStringOption(i, 0)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get username option")
		return
	}

	err = data.AddAccountUseCase.Execute(i.GuildID, i.Member.User.ID, username)
	if err != nil {
		if strings.Contains(err.Error(), "already tracked") {
			handleAccountAlreadyTrackedError(s, i, username)
			return
		}
		if strings.Contains(err.Error(), "does not exist") {
			handlePlayerNotFoundError(s, i, username)
			return
		}
		handleCommandError(s, i, err, "Failed to add account")
		return
	}

	description := fmt.Sprintf("Successfully added account **%s**\n\nYour account is now being tracked for competitions.\nIf there's an active BOTM competition, you've been automatically enrolled!", username)
	embed := createSuccessEmbed("✅ Account Added Successfully", description, i)
	sendEmbedResponse(s, i, embed)
}
