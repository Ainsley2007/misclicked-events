package commands

import (
	"fmt"
	"misclicked-events/data"
	"misclicked-events/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func StartActivityCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !utils.IsAdmin(i) {
		utils.RespondWithError(s, i, fmt.Errorf("you do not have the required permissions to use this command"))
		return
	}

	currentBoss := data.GetCurrentBoss(i.GuildID)
	if len(currentBoss) > 0 {
		response := fmt.Sprintf("An activity has already been selected: \"**%s**\", You need to end this activity before starting a new one.", currentBoss)
		utils.RespondWithMessage(s, i, "%s", response)
		return
	}

	// Get the selected choice
	choice := i.ApplicationCommandData().Options[0].StringValue()
	password := i.ApplicationCommandData().Options[1].StringValue()

	// Defer the response to indicate processing
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Println("Error deferring the response:", err)
		return
	}

	// Perform the long-running operation
	err = data.StartCompetition(i.GuildID, choice, password)
	if err != nil {
		// Edit the deferred response to indicate an error
		errorMessage := "Something went wrong trying to start this activity."
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errorMessage,
		})
		if err != nil {
			fmt.Println("Error editing interaction response:", err)
		}
		return
	}

	updateCategoryChannelName(s, i.GuildID, choice)

	// Edit the deferred response with the final result
	successMessage := fmt.Sprintf(
		"Activity selected: **%s**, now tracking kc for: **%s**",
		choice,
		strings.Join(data.Activities[choice], ", "),
	)
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &successMessage,
	})
	if err != nil {
		fmt.Println("Error editing interaction response:", err)
	}

}

func updateCategoryChannelName(s *discordgo.Session, guildID, currentBoss string) {
	config, err := data.GetBotConfig(guildID)
	if err != nil {
		return
	}

	if config.CategoryChannelID == "" {
		return
	}

	newName := fmt.Sprintf("╔═══BOTM - %s═══╗", currentBoss)
	_, err = s.ChannelEdit(config.CategoryChannelID, &discordgo.ChannelEdit{
		Name: newName,
	})
	if err != nil {
		return
	}
}
