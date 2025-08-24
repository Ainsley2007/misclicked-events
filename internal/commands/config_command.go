package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/domain"

	"github.com/bwmarrin/discordgo"
)

var ConfigCommand = &discordgo.ApplicationCommand{
	Name:        "setup-channels",
	Description: "Setup channels to show competition results",
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
			Description: "Select a channels category (optional)",
			Required:    false,
		},
	},
}

func HandleConfigCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := deferResponse(s, i, "setup-channels"); err != nil {
		return
	}

	if !IsAdmin(i) {
		handleCommandError(s, i, fmt.Errorf("you don't have the required permissions to use this command"), "Permission check failed")
		return
	}

	if err := validateRequiredOptions(i, 2, "setup-channels"); err != nil {
		handleCommandError(s, i, err, "Config command validation failed")
		return
	}

	applicationCommandData := i.ApplicationCommandData().Options
	rankingChannel := applicationCommandData[0].ChannelValue(s)
	hiscoreChannel := applicationCommandData[1].ChannelValue(s)

	if rankingChannel == nil || hiscoreChannel == nil {
		handleCommandError(s, i, fmt.Errorf("failed to retrieve channel information"), "Channel retrieval failed")
		return
	}

	rankingChannelID := rankingChannel.ID
	hiscoreChannelID := hiscoreChannel.ID

	if rankingChannel.GuildID != i.GuildID || hiscoreChannel.GuildID != i.GuildID {
		handleCommandError(s, i, fmt.Errorf("all channels must be in this server"), "Channel validation failed")
		return
	}

	var categoryChannelID string
	if len(applicationCommandData) > 2 {
		categoryChannel := applicationCommandData[2].ChannelValue(s)
		if categoryChannel != nil && categoryChannel.GuildID == i.GuildID {
			categoryChannelID = categoryChannel.ID
		} else if categoryChannel != nil {
			handleCommandError(s, i, fmt.Errorf("category channel must be in this server"), "Category channel validation failed")
			return
		}
	}

	config := &domain.Config{
		RankingChannelID:  rankingChannelID,
		HiscoreChannelID:  hiscoreChannelID,
		CategoryChannelID: categoryChannelID,
	}

	err := data.ConfigRepo.SaveConfig(config, i.GuildID)
	if err != nil {
		handleCommandError(s, i, fmt.Errorf("failed to save configuration: %w", err), "Config save failed")
		return
	}

	description := fmt.Sprintf("**Ranking Channel:** <#%s>\n", rankingChannelID)
	description += fmt.Sprintf("**BOTM Channel:** <#%s>", hiscoreChannelID)

	if categoryChannelID != "" {
		description += fmt.Sprintf("\n**Category Channel:** <#%s>", categoryChannelID)
	}

	description += "\n\nYour competition results will now be displayed in these channels."

	embed := createSuccessEmbed("✅ Configuration Saved", description, i)
	sendEmbedResponse(s, i, embed)
}
