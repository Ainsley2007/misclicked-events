package data

import (
	"misclicked-events/internal/utils"
)

func UpdateConfig(guildID, rankingChannelID, hiscoreChannelID, categoryChannelID string) error {
	botConfig := BotConfig{
		RankingChannelID:  rankingChannelID,
		HiscoreChannelID:  hiscoreChannelID,
		CategoryChannelID: categoryChannelID,
	}

	err := SaveBotConfig(guildID, botConfig)
	if err != nil {
		utils.LogError("Something went wrong while updating config", err)
	}
	return err
}
