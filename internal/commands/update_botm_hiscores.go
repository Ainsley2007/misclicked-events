package commands

import (
	"fmt"
	"misclicked-events/internal/constants"
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"
	"time"

	"github.com/bwmarrin/discordgo"
)

func UpdateBOTMHiscores(s *discordgo.Session) {
	// Do an initial run when the bot starts
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
		} else {
			err := updateNoEventMessage(s, guild.ID)
			if err != nil {
				utils.LogError("Error when updating no-event message", err)
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

	currentActivity := data.GetCurrentBoss(guildID)
	if currentActivity == "" {
		return fmt.Errorf("no event found")
	}

	// Fetch the leaderboard data
	participantKC, err := data.GetParticipantsByActivityKCThreshold(guildID, currentActivity)
	if err != nil {
		return fmt.Errorf("error fetching participants: %w", err)
	}

	// Build the embed
	embed := &discordgo.MessageEmbed{
		Title: "🏆 Killcount Leaderboard",
		Color: 0xffd700, // Gold for leaderboard
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: constants.Activities[currentActivity].BossThumbnail, // Replace with a relevant boss icon
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("🔄 Last updated: %s", time.Now().Format("Jan 02, 2006 15:04:05 MST")),
		},
	}

	// Add information about tracked bosses
	embed.Description = "### Tracked Bosses:\n"
	bosses := constants.Activities[currentActivity].BossNames
	for _, boss := range bosses {
		embed.Description += fmt.Sprintf("• *%s*\n", boss)
	}

	// Add leaderboard details
	if len(participantKC) == 0 {
		embed.Description += "\n🚨 No participants have enough KC yet!\n"
	} else {
		embed.Description += "### Leaderboard:\n"
		// Initialize rank tracking variables
		rank := 0
		previousKC := -1 // Set to a value that cannot match any valid TotalKC
		currentRank := 1 // The rank being assigned to participants

		for _, participant := range participantKC {
			// Check if the current participant's TotalKC is different from the previous one
			if participant.TotalKC != previousKC {
				currentRank = rank + 1 // Update the current rank
			}

			// Assign appropriate emoji for ranks
			var rankEmoji string
			switch currentRank {
			case 1:
				rankEmoji = "🥇" // Gold Medal
			case 2:
				rankEmoji = "🥈" // Silver Medal
			case 3:
				rankEmoji = "🥉" // Bronze Medal
			default:
				rankEmoji = fmt.Sprintf("%d.", currentRank) // Numeric ranking for 4th and beyond
			}

			// Build account-specific details
			accountDetails := ""
			for _, account := range participant.AccountKCs {
				accountDetails += fmt.Sprintf("\u00A0\u00A0\u00A0\u00A0 ┗ *%s: %d*\n", account.AccountName, account.TotalKC)
			}

			// Add the rank, mention, total KC, and account details to the description
			embed.Description += fmt.Sprintf(
				"%s **<@%s>** - **Total KC:** `%d`\n%s\n",
				rankEmoji, participant.DiscordId, participant.TotalKC, accountDetails,
			)

			// Update rank and previousKC
			rank++
			previousKC = participant.TotalKC
		}

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

func updateNoEventMessage(s *discordgo.Session, guildID string) error {
	// Fetch the bot configuration for the guild
	config, err := data.GetBotConfig(guildID)
	if err != nil {
		return fmt.Errorf("error fetching bot configuration: %w", err)
	}

	// Build the embed
	embed := &discordgo.MessageEmbed{
		Title:       "🚨 No Ongoing Event",
		Color:       0xFFA500, // Orange for no event
		Description: "There is currently no ongoing event. Use the appropriate command to start a new event!",
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("🆕 Last updated: %s", time.Now().Format("Jan 02, 2006 15:04:05 MST")),
		},
	}

	// Post or update the no-event message
	if config.HiscoreMessageID != "" {
		// Try to update the existing message
		_, err := s.ChannelMessageEditEmbed(config.HiscoreChannelID, config.HiscoreMessageID, embed)
		if err != nil {
			// If editing fails, post a new message and update the message ID
			newMessage, err := s.ChannelMessageSendEmbed(config.HiscoreChannelID, embed)
			if err != nil {
				return fmt.Errorf("error sending new no-event message: %w", err)
			}
			data.UpdateHiscoreMessageID(guildID, newMessage.ID)
		}
	} else {
		// No previous message, post a new one
		newMessage, err := s.ChannelMessageSendEmbed(config.HiscoreChannelID, embed)
		if err != nil {
			return fmt.Errorf("error sending no-event message: %w", err)
		}
		data.UpdateHiscoreMessageID(guildID, newMessage.ID)
	}

	return nil
}
