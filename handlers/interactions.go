package handlers

import (
	"misclicked-events/commands"
	"misclicked-events/utils"

	"github.com/bwmarrin/discordgo"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "setup-channels":
		commands.ConfigCommand(s, i)
	case "track":
		commands.TrackNewAccountCommand(s, i)
	case "untrack":
		commands.UnTrackAccountCommand(s, i)
	case "tracking":
		commands.TrackedAccountsCommand(s, i)
	case "start":
		commands.StartActivityCommand(s, i)
	case "end":
		commands.EndActivityCommand(s, i)
	default:
		utils.LogError("Unknown command", nil)
	}
}
