package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/domain"
	"misclicked-events/internal/utils"

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
	// Defer the response immediately to prevent timeout
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		utils.Error("Failed to defer response for config command: %v", err)
		return
	}

	// Check permissions
	if !utils.IsAdmin(i) {
		utils.EditResponseError(s, i, fmt.Errorf("you don't have the required permissions to use this command"))
		return
	}

	// Parse command options
	applicationCommandData := i.ApplicationCommandData().Options
	if len(applicationCommandData) < 2 {
		utils.EditResponseError(s, i, fmt.Errorf("invalid command options: please provide both ranking channels"))
		return
	}

	// Extract channel IDs
	rankingChannel := applicationCommandData[0].ChannelValue(s)
	hiscoreChannel := applicationCommandData[1].ChannelValue(s)

	if rankingChannel == nil || hiscoreChannel == nil {
		utils.EditResponseError(s, i, fmt.Errorf("failed to retrieve channel information"))
		return
	}

	rankingChannelID := rankingChannel.ID
	hiscoreChannelID := hiscoreChannel.ID

	// Validate that channels are in the same guild
	if rankingChannel.GuildID != i.GuildID || hiscoreChannel.GuildID != i.GuildID {
		utils.EditResponseError(s, i, fmt.Errorf("all channels must be in this server"))
		return
	}

	// Handle optional category channel
	var categoryChannelID string
	if len(applicationCommandData) > 2 {
		categoryChannel := applicationCommandData[2].ChannelValue(s)
		if categoryChannel != nil && categoryChannel.GuildID == i.GuildID {
			categoryChannelID = categoryChannel.ID
		} else if categoryChannel != nil {
			utils.EditResponseError(s, i, fmt.Errorf("category channel must be in this server"))
			return
		}
	}

	// Create config object
	config := &domain.Config{
		RankingChannelID:  rankingChannelID,
		HiscoreChannelID:  hiscoreChannelID,
		CategoryChannelID: categoryChannelID,
	}

	// Save configuration
	err = data.ConfigRepo.SaveConfig(config, i.GuildID)
	if err != nil {
		utils.EditResponseError(s, i, fmt.Errorf("failed to save configuration: %w", err))
		return
	}

	// Build success message
	successMessage := "✅ **Configuration saved successfully!**\n\n"
	successMessage += fmt.Sprintf("**Ranking Channel:** <#%s>\n", rankingChannelID)
	successMessage += fmt.Sprintf("**BOTM Channel:** <#%s>", hiscoreChannelID)

	if categoryChannelID != "" {
		successMessage += fmt.Sprintf("\n**Category Channel:** <#%s>", categoryChannelID)
	}

	successMessage += "\n\nYour competition results will now be displayed in these channels."

	utils.EditResponseMessage(s, i, successMessage)
}
