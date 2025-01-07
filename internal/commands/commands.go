package commands

import (
	"fmt"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func RegisterCommands(s *discordgo.Session, force bool) {

	commands := []*discordgo.ApplicationCommand{
		ConfigCommand,
		TrackAccountCommand,
		UntrackAccountCommand,
		TrackedAccountsCommand,
		StartActivityCommand,
		EndActivityCommand,
	}

	existingCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		utils.LogError("Error fetching existing commands:", err)
		return
	}

	// Check if re-registration is necessary
	if !force && commandsAreEqual(existingCommands, commands) {
		fmt.Println("Commands are already up-to-date. Skipping registration.")
		return
	}

	// Bulk overwrite commands
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", commands)
	if err != nil {
		utils.LogError("Error overwriting commands:", err)
	} else {
		fmt.Println("Commands registered successfully.")
	}
}

func commandsAreEqual(existing []*discordgo.ApplicationCommand, new []*discordgo.ApplicationCommand) bool {
	if len(existing) != len(new) {
		return false
	}

	// Create a map of existing commands for easier lookup
	existingMap := make(map[string]*discordgo.ApplicationCommand)
	for _, cmd := range existing {
		existingMap[cmd.Name] = cmd
	}

	for _, newCmd := range new {
		// Check if the command exists in the map
		existingCmd, ok := existingMap[newCmd.Name]
		if !ok {
			return false // Command is missing
		}

		// Compare descriptions
		if newCmd.Description != existingCmd.Description {
			return false
		}

		// Compare options
		if !optionsAreEqual(newCmd.Options, existingCmd.Options) {
			return false
		}
	}

	return true
}

func optionsAreEqual(newOpts []*discordgo.ApplicationCommandOption, existingOpts []*discordgo.ApplicationCommandOption) bool {
	if len(newOpts) != len(existingOpts) {
		return false
	}

	// Compare each option
	for i, newOpt := range newOpts {
		existingOpt := existingOpts[i]
		if newOpt.Name != existingOpt.Name ||
			newOpt.Description != existingOpt.Description ||
			newOpt.Type != existingOpt.Type ||
			newOpt.Required != existingOpt.Required {
			return false
		}

		// Compare choices (if applicable)
		if !choicesAreEqual(newOpt.Choices, existingOpt.Choices) {
			return false
		}
	}

	return true
}

func choicesAreEqual(newChoices []*discordgo.ApplicationCommandOptionChoice, existingChoices []*discordgo.ApplicationCommandOptionChoice) bool {
	if len(newChoices) != len(existingChoices) {
		return false
	}

	// Compare each choice
	for i, newChoice := range newChoices {
		existingChoice := existingChoices[i]
		if newChoice.Name != existingChoice.Name || newChoice.Value != existingChoice.Value {
			return false
		}
	}

	return true
}
