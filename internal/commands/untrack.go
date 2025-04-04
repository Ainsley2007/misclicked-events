package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

var UntrackAccountCommand = &discordgo.ApplicationCommand{
	Name:        "untrack",
	Description: "Untracks an OSRS account from your profile",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "username",
			Description: "The OSRS account username to stop tracking",
			Required:    true,
		},
	},
}

func HandleUnTrackAccountCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	// Ensure the required option is provided
	options := i.ApplicationCommandData().Options
	if len(options) < 1 {
		utils.EditResponseError(s, i, fmt.Errorf("please provide the OSRS account username you want to stop tracking"))
		return
	}

	// Extract the username from the command
	username := options[0].StringValue()

	// Attempt to untrack the account
	err = data.UntrackAccount(i.GuildID, username, i.Member.User.ID)
	if err != nil {
		utils.EditResponseError(s, i, fmt.Errorf("could not untrack the account '%s': %w", username, err))
		return
	}

	//update the hiscore message
	err = UpdateHiscoreMessage(s, i.GuildID)
	if err != nil {
		utils.LogError("Error updating hiscore message", err)
	}

	// Respond with success message
	response := fmt.Sprintf("Successfully stopped tracking the OSRS account: **%s**.", username)
	utils.EditResponseMessage(s, i, response)
}
