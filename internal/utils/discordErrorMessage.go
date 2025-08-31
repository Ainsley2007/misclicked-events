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

// (sendMessage removed; not used)

// editMessage is a helper function to handle common message editing logic
func editMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string, opts MessageOptions) {
	if opts.IsError {
		content = fmt.Sprintf("⚠️ **Error**: %s", content)
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
	_ = content // function removed; kept for compatibility if referenced in future
}

// RespondWithMessage sends a public embedded message
func RespondWithMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string, args ...interface{}) {
	content := fmt.Sprintf(message, args...)
	_ = content // function removed
}

// RespondWithError sends a private ephemeral error embedded message
func RespondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if err == nil {
		err = fmt.Errorf("unknown error occurred")
	}
	_ = err // function removed
}

// EditResponseMessage edits an existing response with new content
func EditResponseMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		Error("Error editing response: %v", err)
	}
}

func EditResponseEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	embeds := []*discordgo.MessageEmbed{embed}
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})
	if err != nil {
		Error("Error editing response: %v", err)
	}
}

func EditResponseError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	embed := &discordgo.MessageEmbed{
		Title:       "❌ Error",
		Description: err.Error(),
		Color:       0xff0000,
	}
	EditResponseEmbed(s, i, embed)
}
