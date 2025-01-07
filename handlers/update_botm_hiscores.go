package handlers

import (
	"fmt"
	"misclicked-events/data"
	"misclicked-events/utils"
	"time"

	"github.com/bwmarrin/discordgo"
)

func UpdateBOTMHiscores(s *discordgo.Session) {
	//Do an initial run when the bot starts
	updateUsers(s)

	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			updateUsers(s)
		}
	}
}

func updateUsers(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		if ongoingEvent := checkOngoingEvent(guild.ID); ongoingEvent != "" {
			err := data.UpdateAccountsKC(guild.ID)
			if err != nil {
				utils.LogError("Error when updating accounts", err)
				continue
			}

			err = UpdateHiscoreMessage(s, guild.ID)
			if err != nil {
				utils.LogError("Error when updating hiscore message", err)
			}

		}
	}
}

func checkOngoingEvent(guildID string) string {
	currentEvent := data.GetCurrentBoss(guildID)
	if currentEvent == "" {
		return ""
	}
	return currentEvent
}

func UpdateHiscoreMessage(s *discordgo.Session, guildID string) error {
	// Fetch the bot configuration for the guild
	config, err := data.GetBotConfig(guildID)
	if err != nil {
		return fmt.Errorf("error fetching bot configuration: %w", err)
	}

	// Fetch the leaderboard data
	participantKC, err := data.GetParticipantsByKCThreshold(guildID, 10)
	if err != nil {
		return fmt.Errorf("error fetching participants: %w", err)
	}

	// Build the embed
	embed := &discordgo.MessageEmbed{
		Title: "🏆 Killcount Leaderboard",
		Color: 0xffd700, // Gold for leaderboard
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://example.com/boss-icon.png", // Replace with a relevant boss icon
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("🆕 Last updated: %s", time.Now().Format("Jan 02, 2006 15:04:05 MST")),
		},
	}

	// Populate the embed description with leaderboard data
	for rank, participant := range participantKC {
		// Assign appropriate emoji for ranks
		var rankEmoji string
		switch rank {
		case 0:
			rankEmoji = "🥇" // Gold Medal
		case 1:
			rankEmoji = "🥈" // Silver Medal
		case 2:
			rankEmoji = "🥉" // Bronze Medal
		default:
			rankEmoji = fmt.Sprintf("%d.", rank+1) // Numeric ranking for 4th and beyond
		}

		// Build account-specific details
		accountDetails := ""
		for _, account := range participant.AccountKCs {
			accountDetails += fmt.Sprintf("• **%s**: %d KC\n", account.AccountName, account.TotalKC)
		}

		// Add the rank, mention, total KC, and account details to the description
		embed.Description += fmt.Sprintf(
			"%s **<@%s>** - Total KC: **%d**\n%s\n",
			rankEmoji, participant.DiscordId, participant.TotalKC, accountDetails,
		)
	}

	// Post or update the leaderboard message
	if config.HiscoreMessageID != "" {
		// Try to update the existing message
		_, err := s.ChannelMessageEditEmbed(config.HiscoreChannelID, config.HiscoreMessageID, embed)
		if err != nil {
			// If editing fails, post a new message and update the message ID
			newMessage, err := s.ChannelMessageSendEmbed(config.HiscoreChannelID, embed)
			if err != nil {
				return fmt.Errorf("error sending new leaderboard message: %w", err)
			}
			data.UpdateHiscoreMessageID(guildID, newMessage.ID)
		}
	} else {
		// No previous message, post a new one
		newMessage, err := s.ChannelMessageSendEmbed(config.HiscoreChannelID, embed)
		if err != nil {
			return fmt.Errorf("error sending leaderboard message: %w", err)
		}
		data.UpdateHiscoreMessageID(guildID, newMessage.ID)
	}

	return nil
}
