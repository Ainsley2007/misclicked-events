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
	if !utils.IsAdmin(i) {
		utils.RespondWithError(s, i, fmt.Errorf("you do not have the required permissions to use this command"))
		return
	}

	password := i.ApplicationCommandData().Options[0].StringValue()

	err := data.EndCompetition(i.GuildID, password)

	if err != nil {
		response := fmt.Sprintf("something went wrong trying to end the event: **%s**", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Activity ended",
		},
	})

}
