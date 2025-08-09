package commands

import (
	"fmt"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func deferResponse(s *discordgo.Session, i *discordgo.InteractionCreate, commandName string) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		utils.Error("Failed to defer response for %s command: %v", commandName, err)
		return err
	}
	return nil
}

func validateRequiredOptions(i *discordgo.InteractionCreate, requiredCount int, commandName string) error {
	applicationCommandData := i.ApplicationCommandData().Options
	if len(applicationCommandData) < requiredCount {
		return fmt.Errorf("invalid command options: please provide all required parameters")
	}
	return nil
}

func getStringOption(i *discordgo.InteractionCreate, index int) (string, error) {
	applicationCommandData := i.ApplicationCommandData().Options
	if index >= len(applicationCommandData) {
		return "", fmt.Errorf("option at index %d not found", index)
	}

	value := applicationCommandData[index].StringValue()
	if value == "" {
		return "", fmt.Errorf("option value cannot be empty")
	}

	return value, nil
}

func handleCommandError(s *discordgo.Session, i *discordgo.InteractionCreate, err error, context string) {
	utils.Error("%s: %v", context, err)
	utils.EditResponseError(s, i, err)
}
