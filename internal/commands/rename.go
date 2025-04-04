package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/service"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

var RenameAccountCommand = &discordgo.ApplicationCommand{
	Name:        "rename",
	Description: "Rename one of your tracked OSRS accounts",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "old_username",
			Description: "The current username of the account",
			Required:    true,
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
	options := i.ApplicationCommandData().Options
	if len(options) < 2 {
		utils.RespondWithError(s, i, fmt.Errorf("please provide both the old and new usernames"))
		return
	}

	oldUsername := options[0].StringValue()
	newUsername := options[1].StringValue()

	// Verify the new username exists in OSRS
	if !service.CheckIfPlayerExists(newUsername) {
		utils.RespondWithError(s, i, fmt.Errorf("could not find an OSRS account with the username: %s", newUsername))
		return
	}

	err := data.RenameAccount(i.GuildID, oldUsername, newUsername, i.Member.User.ID)
	if err != nil {
		utils.RespondWithError(s, i, fmt.Errorf("could not rename the account: %w", err))
		return
	}

	// Update the hiscore message if there's an ongoing event
	if ongoingEvent := checkOngoingEvent(i.GuildID); ongoingEvent != "" {
		err = UpdateHiscoreMessage(s, i.GuildID)
		if err != nil {
			utils.LogError("Error updating hiscore message", err)
		}
	}

	response := fmt.Sprintf("Successfully renamed account from **%s** to **%s**", oldUsername, newUsername)
	utils.RespondWithPrivateMessage(s, i, "%s", response)
}
