package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

var TrackAccountCommand = &discordgo.ApplicationCommand{
	Name:        "track",
	Description: "Link an OSRS account to your profile to track its progress — only add accounts you own.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "username",
			Description: "The username to start tracking",
			Required:    true,
		},
	},
}

func HandleTrackNewAccountCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) < 1 {
		utils.RespondWithError(s, i, fmt.Errorf("please provide a valid username"))
		return
	}

	username := options[0].StringValue()
	success := data.TrackAccount(i.GuildID, username, i.Member.User.ID)
	if !success {
		utils.RespondWithError(s, i, fmt.Errorf("something went wrong trying to track your account"))
		return
	}

	response := fmt.Sprintf("The OSRS account '%s' is now linked to your profile and being tracked for you!", username)
	utils.RespondWithPrivateMessage(s, i, "%s", response)
}
