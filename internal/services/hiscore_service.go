package services

import (
	"fmt"
	"misclicked-events/internal/domain/usecase"
	"misclicked-events/internal/utils"
	"time"

	"github.com/bwmarrin/discordgo"
)

// HiscoreService handles the background task of updating hiscores and tracked accounts
type HiscoreService struct {
	session          *discordgo.Session
	ticker           *time.Ticker
	stop             chan struct{}
	updateHiscoresUC *usecase.UpdateHiscoresUseCase
	done             chan struct{}
}

// NewHiscoreService creates a new hiscore service
func NewHiscoreService(session *discordgo.Session, updateHiscoresUC *usecase.UpdateHiscoresUseCase) *HiscoreService {
	return &HiscoreService{
		session:          session,
		stop:             make(chan struct{}),
		updateHiscoresUC: updateHiscoresUC,
		done:             make(chan struct{}),
	}
}

// Start begins the background task of updating hiscores every hour
func (hs *HiscoreService) Start() {
	// Initial update when bot starts
	hs.updateAllGuilds()

	// Set up ticker for hourly updates
	hs.ticker = time.NewTicker(60 * time.Minute)

	go func() {
		defer close(hs.done)
		defer hs.ticker.Stop()

		for {
			select {
			case <-hs.ticker.C:
				hs.updateAllGuilds()
			case <-hs.stop:
				return
			}
		}
	}()
}

// Stop stops the background task and waits for it to finish
func (hs *HiscoreService) Stop() {
	if hs.ticker != nil {
		hs.ticker.Stop()
	}
	close(hs.stop)
	<-hs.done // Wait for the goroutine to finish
}

// IsRunning returns true if the service is currently running
func (hs *HiscoreService) IsRunning() bool {
	select {
	case <-hs.done:
		return false
	default:
		return true
	}
}

// updateAllGuilds updates all guilds with ongoing events
func (hs *HiscoreService) updateAllGuilds() {
	for _, guild := range hs.session.State.Guilds {
		// Execute the usecase to get all the data we need
		result, err := hs.updateHiscoresUC.Execute(guild.ID)
		if err != nil {
			utils.LogError(fmt.Sprintf("Error when updating hiscores for guild %s", guild.ID), err)
			continue
		}

		// Update Discord messages based on the result
		if result.HasOngoingEvent {
			err = hs.updateHiscoreMessage(result, guild.ID)
			if err != nil {
				utils.LogError(fmt.Sprintf("Error when updating hiscore message for guild %s", guild.ID), err)
			}
		} else {
			err = hs.updateNoEventMessage(result, guild.ID)
			if err != nil {
				utils.LogError(fmt.Sprintf("Error when updating no-event message for guild %s", guild.ID), err)
			}
		}
	}
}

// updateHiscoreMessage updates or creates the hiscore message for a guild
func (hs *HiscoreService) updateHiscoreMessage(result *usecase.UpdateHiscoresResult, guildID string) error {
	// Build the embed
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("🏆 %s Leaderboard", result.EventName),
		Color: 0xffd700, // Gold for leaderboard
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("🔄 Last updated: %s", result.LastUpdated.Format("Jan 02, 2006 15:04:05 MST")),
		},
	}

	// Add information about tracked bosses
	embed.Description = "### Tracked Bosses:\n"
	for _, boss := range result.BossNames {
		embed.Description += fmt.Sprintf("• *%s*\n", boss)
	}

	// Add leaderboard details
	if len(result.Participants) == 0 {
		embed.Description += "\n🚨 No participants have enough KC yet!\n"
	} else {
		embed.Description += "### Leaderboard:\n"
		// Initialize rank tracking variables
		rank := 0
		previousKC := -1 // Set to a value that cannot match any valid KC
		currentRank := 1 // The rank being assigned to participants

		for _, participant := range result.Participants {
			// Calculate total KC for this participant
			totalKC := 0
			for _, account := range participant.Accounts {
				totalKC += account.KCGained
			}

			// Check if the current participant's total KC is different from the previous one
			if totalKC != previousKC {
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
			for _, account := range participant.Accounts {
				accountDetails += fmt.Sprintf("\u00A0\u00A0\u00A0\u00A0 ┗ *%s: %d*\n", account.Username, account.KCGained)
			}

			// Add the rank, mention, total KC, and account details to the description
			embed.Description += fmt.Sprintf(
				"%s **<@%s>** - **Total KC:** `%d`\n%s\n",
				rankEmoji, participant.DiscordID, totalKC, accountDetails,
			)

			// Update rank and previousKC
			rank++
			previousKC = totalKC
		}

		embed.Description += fmt.Sprintf("_Threshold: %dkc_\n", result.Threshold)
	}

	// Post or update the leaderboard message
	if result.Config.HiscoreMessageID != "" {
		// Try to update the existing message
		_, err := hs.session.ChannelMessageEditEmbed(result.Config.HiscoreChannelID, result.Config.HiscoreMessageID, embed)
		if err != nil {
			// If editing fails, post a new message and update the message ID
			newMessage, err := hs.session.ChannelMessageSendEmbed(result.Config.HiscoreChannelID, embed)
			if err != nil {
				return fmt.Errorf("error sending new leaderboard message: %w", err)
			}
			hs.updateHiscoresUC.UpdateHiscoreMessageID(guildID, newMessage.ID)
		}
	} else {
		// No previous message, post a new one
		newMessage, err := hs.session.ChannelMessageSendEmbed(result.Config.HiscoreChannelID, embed)
		if err != nil {
			return fmt.Errorf("error sending leaderboard message: %w", err)
		}
		hs.updateHiscoresUC.UpdateHiscoreMessageID(guildID, newMessage.ID)
	}

	return nil
}

// updateNoEventMessage updates or creates the no-event message for a guild
func (hs *HiscoreService) updateNoEventMessage(result *usecase.UpdateHiscoresResult, guildID string) error {
	// Build the embed
	embed := &discordgo.MessageEmbed{
		Title:       "🚨 No Ongoing Event",
		Color:       0xFFA500, // Orange for no event
		Description: "There is currently no ongoing event. Use the appropriate command to start a new event!",
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("🆕 Last updated: %s", result.LastUpdated.Format("Jan 02, 2006 15:04:05 MST")),
		},
	}

	// Post or update the no-event message
	if result.Config.HiscoreMessageID != "" {
		// Try to update the existing message
		_, err := hs.session.ChannelMessageEditEmbed(result.Config.HiscoreChannelID, result.Config.HiscoreMessageID, embed)
		if err != nil {
			// If editing fails, post a new message and update the message ID
			newMessage, err := hs.session.ChannelMessageSendEmbed(result.Config.HiscoreChannelID, embed)
			if err != nil {
				return fmt.Errorf("error sending new no-event message: %w", err)
			}
			hs.updateHiscoresUC.UpdateHiscoreMessageID(guildID, newMessage.ID)
		}
	} else {
		// No previous message, post a new one
		newMessage, err := hs.session.ChannelMessageSendEmbed(result.Config.HiscoreChannelID, embed)
		if err != nil {
			return fmt.Errorf("error sending no-event message: %w", err)
		}
		hs.updateHiscoresUC.UpdateHiscoreMessageID(guildID, newMessage.ID)
	}

	return nil
}
