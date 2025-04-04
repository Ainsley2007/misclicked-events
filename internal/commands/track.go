package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

var TrackAccountCommand = &discordgo.ApplicationCommand{
	Name:        "track",
	Description: "Link an OSRS account to your profile to track its progress — only add accounts you own.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "username",
			Description: "The username to start tracking",
			Required:    true,
		},
	},
}

func HandleTrackNewAccountCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	options := i.ApplicationCommandData().Options
	if len(options) < 1 {
		utils.EditResponseError(s, i, fmt.Errorf("please provide a valid username"))
		return
	}

	username := options[0].StringValue()
	err = data.TrackAccount(i.GuildID, username, i.Member.User.ID)
	if err != nil {
		utils.EditResponseError(s, i, fmt.Errorf("could not track the account '%s': %w", username, err))
		return
	}

	response := fmt.Sprintf("Successfully started tracking the OSRS account: **%s**", username)
	utils.EditResponseMessage(s, i, response)
}
