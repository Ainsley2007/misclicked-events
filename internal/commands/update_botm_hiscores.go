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
	config, err := data.ConfigRepo.FetchConfig(guildID)
	if err != nil {
		return fmt.Errorf("error fetching bot configuration: %w", err)
	}

	currentActivity := data.GetCurrentBoss(guildID)
	if currentActivity == "" {
		return fmt.Errorf("no event found")
	}

	participantKC, err := data.GetParticipantsByActivityKCThreshold(guildID)
	if err != nil {
		return fmt.Errorf("error fetching participants: %w", err)
	}

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

	embed.Description = "### Tracked Bosses:\n"
	bosses := constants.Activities[currentActivity].BossNames
	for _, boss := range bosses {
		embed.Description += fmt.Sprintf("• *%s*\n", boss)
	}

	if len(participantKC) == 0 {
		embed.Description += "\n🚨 No participants have enough KC yet!\n"
	} else {
		embed.Description += "### Leaderboard:\n"
		rank := 0
		previousKC := -1
		currentRank := 1

		for _, participant := range participantKC {
			if participant.TotalKC != previousKC {
				currentRank = rank + 1
			}

			var rankEmoji string
			switch currentRank {
			case 1:
				rankEmoji = "🥇"
			case 2:
				rankEmoji = "🥈"
			case 3:
				rankEmoji = "🥉"
			default:
				rankEmoji = fmt.Sprintf("%d.", currentRank)
			}

			accountDetails := ""
			for _, account := range participant.AccountKCs {
				accountDetails += fmt.Sprintf("\u00A0\u00A0\u00A0\u00A0 ┗ *%s: %d*\n", account.AccountName, account.TotalKC)
			}

			embed.Description += fmt.Sprintf(
				"%s **<@%s>** - **Total KC:** `%d`\n%s\n",
				rankEmoji, participant.DiscordId, participant.TotalKC, accountDetails,
			)

			rank++
			previousKC = participant.TotalKC
		}

		embed.Description += fmt.Sprintf("_Threshold: %dkc_\n", constants.Activities[currentActivity].Threshold)
	}

	if config.HiscoreMessageID != "" {
		_, err := s.ChannelMessageEditEmbed(config.HiscoreChannelID, config.HiscoreMessageID, embed)
		if err != nil {
			newMessage, err := s.ChannelMessageSendEmbed(config.HiscoreChannelID, embed)
			if err != nil {
				return fmt.Errorf("error sending new leaderboard message: %w", err)
			}
			data.ConfigRepo.EditHiscoreMessageID(guildID, newMessage.ID)
		}
	} else {
		newMessage, err := s.ChannelMessageSendEmbed(config.HiscoreChannelID, embed)
		if err != nil {
			return fmt.Errorf("error sending leaderboard message: %w", err)
		}
		data.ConfigRepo.EditHiscoreMessageID(guildID, newMessage.ID)
	}

	return nil
}

func updateNoEventMessage(s *discordgo.Session, guildID string) error {
	config, err := data.ConfigRepo.FetchConfig(guildID)
	if err != nil {
		return fmt.Errorf("error fetching bot configuration: %w", err)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🚨 No Ongoing Event",
		Color:       0xFFA500,
		Description: "There is currently no ongoing event. Use the appropriate command to start a new event!",
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("🆕 Last updated: %s", time.Now().Format("Jan 02, 2006 15:04:05 MST")),
		},
	}

	if config.HiscoreMessageID != "" {
		_, err := s.ChannelMessageEditEmbed(config.HiscoreChannelID, config.HiscoreMessageID, embed)
		if err != nil {
			newMessage, err := s.ChannelMessageSendEmbed(config.HiscoreChannelID, embed)
			if err != nil {
				return fmt.Errorf("error sending new no-event message: %w", err)
			}
			data.ConfigRepo.EditHiscoreMessageID(guildID, newMessage.ID)
		}
	} else {
		newMessage, err := s.ChannelMessageSendEmbed(config.HiscoreChannelID, embed)
		if err != nil {
			return fmt.Errorf("error sending no-event message: %w", err)
		}
		data.ConfigRepo.EditHiscoreMessageID(guildID, newMessage.ID)
	}

	return nil
}
