package commands

import (
	"fmt"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func RegisterCommands(s *discordgo.Session, force bool) {

	commands := []*discordgo.ApplicationCommand{
		ConfigCommand,
		AddAccountCommand,
		RemoveAccountCommand,
		TrackedAccountsCommand,
		StartActivityCommand,
		EndActivityCommand,
		RenameAccountCommand,
	}

	existingCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		utils.LogError("Error fetching existing commands:", err)
		return
	}

	if !force && commandsAreEqual(existingCommands, commands) {
		fmt.Println("Commands are already up-to-date. Skipping registration.")
		return
	}

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

	existingMap := make(map[string]*discordgo.ApplicationCommand)
	for _, cmd := range existing {
		existingMap[cmd.Name] = cmd
	}

	for _, newCmd := range new {
		existingCmd, ok := existingMap[newCmd.Name]
		if !ok {
			return false
		}

		if newCmd.Description != existingCmd.Description {
			return false
		}

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

	for i, newOpt := range newOpts {
		existingOpt := existingOpts[i]
		if newOpt.Name != existingOpt.Name ||
			newOpt.Description != existingOpt.Description ||
			newOpt.Type != existingOpt.Type ||
			newOpt.Required != existingOpt.Required {
			return false
		}

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

	for i, newChoice := range newChoices {
		existingChoice := existingChoices[i]
		if newChoice.Name != existingChoice.Name || newChoice.Value != existingChoice.Value {
			return false
		}
	}

	return true
}
