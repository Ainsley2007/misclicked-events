package handlers

import (
	"misclicked-events/internal/commands"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		handleApplicationCommand(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		handleAutocomplete(s, i)
	default:
		utils.LogError("Unknown interaction type", nil)
	}
}

func handleApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "setup-channels":
		commands.HandleConfigCommand(s, i)
	case "add-account":
		commands.HandleAddAccountCommand(s, i)
	case "remove-account":
		commands.HandleRemoveAccountCommand(s, i)
	case "tracked-accounts":
		commands.HandleTrackedAccountsCommand(s, i)
	case "start":
		commands.HandleStartActivityCommand(s, i)
	case "end":
		commands.HandleEndActivityCommand(s, i)
	case "rename-account":
		commands.HandleRenameAccountCommand(s, i)
	default:
		utils.LogError("Unknown command", nil)
	}
}

func handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "remove-account":
		commands.HandleAccountAutocomplete(s, i)
	case "rename-account":
		commands.HandleAccountAutocomplete(s, i)
	default:
		utils.LogError("Unknown autocomplete command", nil)
	}
}
