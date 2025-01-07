package commands

import (
	"fmt"
	"misclicked-events/data"

	"github.com/bwmarrin/discordgo"
)

func UnTrackAccountCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) < 1 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a username.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	username := options[0].StringValue()
	err := data.UntrackAccount(i.GuildID, username, i.Member.User.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error untracking account: %s", err),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	response := fmt.Sprintf("Stopped tracking: %s", username)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
