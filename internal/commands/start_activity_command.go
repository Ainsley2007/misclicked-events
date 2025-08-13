package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var StartActivityCommand = &discordgo.ApplicationCommand{
	Name:        "start",
	Description: "Select an activity to start",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "choice",
			Description:  "Choose an activity",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "password",
			Description: "Set an activity password",
			Required:    true,
		},
	},
}

func HandleStartActivityCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !IsAdmin(i) {
		utils.RespondWithError(s, i, fmt.Errorf("you do not have the required permissions to use this command"))
		return
	}

	choice := i.ApplicationCommandData().Options[0].StringValue()
	password := i.ApplicationCommandData().Options[1].StringValue()

	err := deferResponse(s, i, "start-activity")
	if err != nil {
		return
	}

	botm, err := data.StartActivityUseCase.Execute(i.GuildID, choice, password)
	if err != nil {
		handleCommandError(s, i, err, "Failed to start activity")
		return
	}

	updateCategoryChannelName(s, i.GuildID, choice)

	successMessage := fmt.Sprintf(
		"Activity selected: **%s**, now tracking kc for: **%s**",
		choice,
		strings.Join(botm.Activity.HiscoreNames, ", "),
	)
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &successMessage,
	})
	if err != nil {
		fmt.Println("Error editing interaction response:", err)
	}
}

func updateCategoryChannelName(s *discordgo.Session, guildID, currentBoss string) {
	config, err := data.ConfigRepo.FetchConfig(guildID)
	if err != nil {
		fmt.Println("error fetching bot configuration: %w", err)
		return
	}

	if config.CategoryChannelID == "" {
		return
	}

	newName := fmt.Sprintf("╔═══BOTM - %s═══╗", currentBoss)
	_, err = s.ChannelEdit(config.CategoryChannelID, &discordgo.ChannelEdit{
		Name: newName,
	})
	if err != nil {
		return
	}
}
