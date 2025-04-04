package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type MessageOptions struct {
	IsEphemeral bool
	IsError     bool
	Color       int
}

// sendMessage is a helper function to handle common message sending logic
func sendMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string, opts MessageOptions) {
	var color int
	if opts.IsError {
		color = 0xff0000 // Red for errors
		content = fmt.Sprintf("⚠️ **Error**\n%s", content)
	} else {
		if opts.Color != 0 {
			color = opts.Color
		} else {
			color = 0x00ccff // Default light blue
		}
	}

	embed := &discordgo.MessageEmbed{
		Description: content,
		Color:       color,
	}

	if opts.IsError {
		embed.Title = "An Error Occurred"
	}

	data := &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{embed},
	}

	if opts.IsEphemeral {
		data.Flags = discordgo.MessageFlagsEphemeral
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
	if err != nil {
		LogError("Failed to send message", err)
	}
}

// editMessage is a helper function to handle common message editing logic
func editMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string, opts MessageOptions) {
	if opts.IsError {
		content = fmt.Sprintf("⚠️ **Error**\n%s", content)
	}

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		LogError("Error editing response", err)
	}
}

// RespondWithPrivateMessage sends a private ephemeral embedded message
func RespondWithPrivateMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string, args ...interface{}) {
	content := fmt.Sprintf(message, args...)
	sendMessage(s, i, content, MessageOptions{
		IsEphemeral: true,
		Color:       0x00ccff,
	})
}

// RespondWithMessage sends a public embedded message
func RespondWithMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string, args ...interface{}) {
	content := fmt.Sprintf(message, args...)
	sendMessage(s, i, content, MessageOptions{
		IsEphemeral: false,
		Color:       0x33cc33,
	})
}

// RespondWithError sends a private ephemeral error embedded message
func RespondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if err == nil {
		err = fmt.Errorf("unknown error occurred")
	}
	sendMessage(s, i, err.Error(), MessageOptions{
		IsEphemeral: true,
		IsError:     true,
	})
}

// EditResponseMessage edits an existing response with new content
func EditResponseMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	editMessage(s, i, content, MessageOptions{})
}

// EditResponseError edits an existing response with an error message
func EditResponseError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if err == nil {
		err = fmt.Errorf("unknown error occurred")
	}
	editMessage(s, i, err.Error(), MessageOptions{
		IsError: true,
	})
}
