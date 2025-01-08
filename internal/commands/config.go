package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

var ConfigCommand = &discordgo.ApplicationCommand{
	Name:        "setup-channels",
	Description: "setup channels to show competition results",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "overall_ranking_channel",
			Description: "Select a channel to show overall competition ranking",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "botm_ranking_channel",
			Description: "Select a channel to show BOTM ranking",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "category_channel",
			Description: "Select a channels category",
			Required:    false,
		},
	},
}

func HandleConfigCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !utils.IsAdmin(i) {
		utils.RespondWithError(s, i, fmt.Errorf("you don't have the required permissions"))
		return
	}

	rankingChannelID := i.ApplicationCommandData().Options[0].ChannelValue(s)
	hiscoreChannelID := i.ApplicationCommandData().Options[1].ChannelValue(s)
	categoryChannelID := i.ApplicationCommandData().Options[2].ChannelValue(s)

	err := data.UpdateConfig(i.GuildID, rankingChannelID.ID, hiscoreChannelID.ID, categoryChannelID.ID)
	if err != nil {
		utils.RespondWithError(s, i, fmt.Errorf("something went wrong while trying to update the config"))
	}

	utils.RespondWithPrivateMessage(s, i, "%s", "Config saved!")
}