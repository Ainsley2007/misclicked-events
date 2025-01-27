package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

var EndActivityCommand = &discordgo.ApplicationCommand{
	Name:        "end",
	Description: "End the current activity",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "password",
			Description: "provide the activity password",
			Required:    true,
		},
	},
}

func HandleEndActivityCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Acknowledge the interaction immediately to prevent timeout
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Printf("Error acknowledging interaction: %v\n", err)
		return
	}

	if !utils.IsAdmin(i) {
		// Send an error message after deferring
		utils.EditResponseMessage(s, i, "❌ You do not have the required permissions to use this command.")
		return
	}

	password := i.ApplicationCommandData().Options[0].StringValue()

	// End the competition
	err = data.EndCompetition(i.GuildID, password)
	if err != nil {
		// Edit the response with an error message
		utils.EditResponseMessage(s, i, "❌ Something went wrong while trying to end the event.")
		return
	}

	// Update the ranking message
	err = updateRankingMessage(s, i.GuildID)
	if err != nil {
		// Edit the response with an error message
		utils.EditResponseMessage(s, i, "❌ Something went wrong while updating the ranking message.")
		return
	}

	// Edit the response to indicate success
	utils.EditResponseMessage(s, i, "✅ The event has ended, and the rankings have been updated!")
}

func updateRankingMessage(s *discordgo.Session, guildID string) error {
	config, err := data.GetBotConfig(guildID)
	if err != nil {
		return fmt.Errorf("error fetching bot configuration: %w", err)
	}

	participants, err := data.GetParticipantsInOrder(guildID)
	if err != nil {
		return fmt.Errorf("error fetching participants: %w", err)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Competition Ranking\n",
		Color:       0x999999,
		Description: "",
	}

	if len(participants) == 0 {
		embed.Description = "🚨 Participants don't have any points yet!"
		return nil
	}

	rank := 0
	previousPoints := -1

	for i, participant := range participants {
		if participant.Points == 0 {
			continue
		}

		var rankEmoji string
		if participant.Points != previousPoints {
			embed.Description += "\n"
			rank = i + 1
			switch rank {
			case 1:
				rankEmoji = "👑" // Victor
			default:
				rankEmoji = fmt.Sprintf("  %d ", rank) // Numeric ranking for 4th and beyond
			}
		} else {
			rankEmoji = "     "
		}

		// Add the participant's rank, username, and points to the embed description
		embed.Description += fmt.Sprintf("%s  **<@%s>**  -  _%d pts_\n", rankEmoji, participant.DiscordId, participant.Points)

		// Update tracking variables
		previousPoints = participant.Points
	}

	// Send or edit the ranking message
	if config.RankingMessageID != "" {
		_, err := s.ChannelMessageEditEmbed(config.RankingChannelID, config.RankingMessageID, embed)
		if err != nil {
			newMessage, err := s.ChannelMessageSendEmbed(config.RankingChannelID, embed)
			if err != nil {
				return fmt.Errorf("error sending new ranking message: %w", err)
			}

			data.UpdateRankingMessageID(guildID, newMessage.ID)
		}
	} else {
		newMessage, err := s.ChannelMessageSendEmbed(config.RankingChannelID, embed)
		if err != nil {
			return fmt.Errorf("error sending ranking message: %w", err)
		}

		data.UpdateRankingMessageID(guildID, newMessage.ID)
	}

	return nil
}
