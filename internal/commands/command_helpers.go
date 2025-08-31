package commands

import (
	"fmt"
	"misclicked-events/internal/utils"
	"strings"
	"time"

	"misclicked-events/internal/data"

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
		return "", fmt.Errorf("option value cannot be empty")
	}

	value := applicationCommandData[index].StringValue()
	if value == "" {
		return "", fmt.Errorf("option value cannot be empty")
	}

	return value, nil
}

func getIntegerOption(i *discordgo.InteractionCreate, index int) (int, error) {
	applicationCommandData := i.ApplicationCommandData().Options
	if index >= len(applicationCommandData) {
		return 0, fmt.Errorf("option at index %d not found", index)
	}

	value := applicationCommandData[index].IntValue()
	return int(value), nil
}

func handleCommandError(s *discordgo.Session, i *discordgo.InteractionCreate, err error, context string) {
	utils.Error("%s: %v", context, err)
	utils.EditResponseError(s, i, err)
}

func HandleAccountAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	accounts, err := data.ParticipantRepo.GetTrackedAccounts(i.GuildID, i.Member.User.ID)
	if err != nil {
		utils.Error("Failed to get tracked accounts for autocomplete: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
		return
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, account := range accounts {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  account,
			Value: account,
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func HandleActivityAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if data.ActivityRepo == nil {
		utils.Error("ActivityRepo is nil")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
		return
	}

	activities, err := data.ActivityRepo.GetAllActivities()
	if err != nil {
		utils.Error("Failed to get activities for autocomplete: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
		return
	}

	focusedOption := i.ApplicationCommandData().Options[0]
	userInput := strings.ToLower(focusedOption.StringValue())

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, activity := range activities {
		if len(choices) >= 25 {
			break
		}

		activityName := strings.ToLower(activity.Name)
		if strings.Contains(activityName, userInput) {
			choiceName := fmt.Sprintf("%s (%s)", activity.Name, activity.Type)
			if len(activity.HiscoreNames) > 1 {
				choiceName = fmt.Sprintf("%s (%s) - %s", activity.Name, activity.Type, strings.Join(activity.HiscoreNames, ", "))
			}

			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  choiceName,
				Value: activity.Name,
			})
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func HandleStartActivityAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if data.ActivityRepo == nil {
		utils.Error("ActivityRepo is nil")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
		return
	}

	activities, err := data.ActivityRepo.GetActivityListByType("boss")
	if err != nil {
		utils.Error("Failed to get boss activities for autocomplete: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
		return
	}

	focusedOption := i.ApplicationCommandData().Options[0]
	userInput := strings.ToLower(focusedOption.StringValue())

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, activity := range activities {
		if len(choices) >= 25 {
			break
		}

		activityName := strings.ToLower(activity.Name)
		if strings.Contains(activityName, userInput) {
			choiceName := fmt.Sprintf("%s - %s", activity.Name, strings.Join(activity.HiscoreNames, ", "))

			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  choiceName,
				Value: activity.Name,
			})
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func IsAdmin(i *discordgo.InteractionCreate) bool {
	return utils.IsAdmin(i)
}

func createSuccessEmbed(title, description string, i *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0x00ff00,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", i.Member.User.Username),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func createErrorEmbed(title, description string, i *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0xff0000,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", i.Member.User.Username),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func createInfoEmbed(title, description string, i *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0x0099ff,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", i.Member.User.Username),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func sendEmbedResponse(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	utils.EditResponseEmbed(s, i, embed)
}

func sendTextResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	utils.EditResponseMessage(s, i, content)
}

func handleAccountNotFoundError(s *discordgo.Session, i *discordgo.InteractionCreate, username string) {
	embed := createErrorEmbed("❌ Account Not Found", "This account is not being tracked.", i)
	sendEmbedResponse(s, i, embed)
}

func handlePlayerNotFoundError(s *discordgo.Session, i *discordgo.InteractionCreate, username string) {
	embed := createErrorEmbed("❌ Player Not Found", "This player does not exist in OSRS.", i)
	sendEmbedResponse(s, i, embed)
}

func handleAccountAlreadyTrackedError(s *discordgo.Session, i *discordgo.InteractionCreate, username string) {
	embed := createErrorEmbed("❌ Account Already Tracked", "This account is already being tracked.", i)
	sendEmbedResponse(s, i, embed)
}
