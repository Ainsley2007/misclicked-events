package commands

import (
	"fmt"
	"misclicked-events/internal/data"
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
		handleCommandError(s, i, fmt.Errorf("you do not have the required permissions to use this command"), "Permission check failed")
		return
	}

	if err := deferResponse(s, i, "start-activity"); err != nil {
		return
	}

	if err := validateRequiredOptions(i, 2, "start-activity"); err != nil {
		handleCommandError(s, i, err, "Start activity command validation failed")
		return
	}

	choice, err := getStringOption(i, 0)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get choice option")
		return
	}

	password, err := getStringOption(i, 1)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get password option")
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
	sendTextResponse(s, i, successMessage)
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
