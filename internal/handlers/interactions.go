package handlers

import (
	"misclicked-events/internal/commands"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "setup-channels":
		commands.HandleConfigCommand(s, i)
	case "track":
		commands.HandleTrackNewAccountCommand(s, i)
	case "untrack":
		commands.HandleUnTrackAccountCommand(s, i)
	case "tracking":
		commands.HandleTrackedAccountsCommand(s, i)
	case "start":
		commands.HandleStartActivityCommand(s, i)
	case "end":
		commands.HandleEndActivityCommand(s, i)
	default:
		utils.LogError("Unknown command", nil)
	}
}
