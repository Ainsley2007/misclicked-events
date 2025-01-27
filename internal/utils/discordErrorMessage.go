package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// RespondWithPrivateMessage sends a private ephemeral embedded message to the user.
func RespondWithPrivateMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string, args ...interface{}) {
	content := fmt.Sprintf(message, args...)
	embed := &discordgo.MessageEmbed{
		Description: content,
		Color:       0x00ccff, // Light blue for private messages
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		LogError("Failed to send private embedded message", err)
	}
}

func EditResponseMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		fmt.Printf("Error editing response message: %v\n", err)
	}
}

// RespondWithMessage sends a public embedded message to the interaction channel.
func RespondWithMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string, args ...interface{}) {
	content := fmt.Sprintf(message, args...)
	embed := &discordgo.MessageEmbed{
		Description: content,
		Color:       0x33cc33, // Green for public messages
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		LogError("Failed to send public embedded message", err)
	}
}

// RespondWithError sends a private ephemeral error embedded message to the user and logs the error.
func RespondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if err == nil {
		err = fmt.Errorf("unknown error occurred")
	}
	content := fmt.Sprintf("⚠️ **Error**\n%s", err.Error())
	embed := &discordgo.MessageEmbed{
		Title:       "An Error Occurred",
		Description: content,
		Color:       0xff0000, // Red for errors
	}

	errSend := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if errSend != nil {
		LogError("Failed to send error embedded message", errSend)
	}
}
