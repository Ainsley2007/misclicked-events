package commands

import (
	"fmt"
	"misclicked-events/internal/data"

	"github.com/bwmarrin/discordgo"
)

var TrackAccountCommand = &discordgo.ApplicationCommand{
	Name:        "track",
	Description: "Tracks killcount for this account and adds it to your discord user.",
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
	err := data.TrackAccount(i.GuildID, username, i.Member.User.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error tracking account: %s", err),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	response := fmt.Sprintf("Now tracking: %s", username)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

}
