package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func RespondWithPrivateMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string, args ...interface{}) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(message, args...),
			Flags:   discordgo.MessageFlagsEphemeral, // Ensure the response is private
		},
	})
}

func RespondWithMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string, args ...interface{}) {
	var content string

	// Check if formatting is needed
	if len(args) > 0 {
		content = fmt.Sprintf(message, args...)
	} else {
		content = message
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral, // Ensure the response is private
		},
	})
}

func RespondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	RespondWithPrivateMessage(s, i, "⚠️ Error: %s", err.Error())
}
