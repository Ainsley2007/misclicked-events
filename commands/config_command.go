package commands

import (
	"fmt"
	"misclicked-events/data"
	"misclicked-events/utils"

	"github.com/bwmarrin/discordgo"
)

func ConfigCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !utils.IsAdmin(i) {
		utils.RespondWithError(s, i, fmt.Errorf("you do not have the required permissions to use this command"))
		return
	}

	rankingChannelID := i.ApplicationCommandData().Options[0].ChannelValue(s)
	hiscoreChannelID := i.ApplicationCommandData().Options[1].ChannelValue(s)
	categoryChannelID := i.ApplicationCommandData().Options[2].ChannelValue(s)

	var config = data.BotConfig{
		CategoryChannelID: categoryChannelID.ID,
		HiscoreChannelID:  hiscoreChannelID.ID,
		RankingChannelID:  rankingChannelID.ID,
	}

	err := data.UpdateChannelIDs(i.GuildID, config)
	if err != nil {
		utils.RespondWithError(s, i, err)
	}

	utils.RespondWithPrivateMessage(s, i, "%s", "Config saved.")
}
