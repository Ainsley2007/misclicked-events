package commands

import (
	"fmt"
	"misclicked-events/internal/data"

	"github.com/bwmarrin/discordgo"
)

var RemoveActivityCommand = &discordgo.ApplicationCommand{
	Name:        "remove-activity",
	Description: "Remove an activity from the system",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "name",
			Description:  "The name of the activity to remove",
			Required:     true,
			Autocomplete: true,
		},
	},
}

func HandleRemoveActivityCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := deferResponse(s, i, "remove-activity"); err != nil {
		return
	}

	if err := validateRequiredOptions(i, 1, "remove-activity"); err != nil {
		handleCommandError(s, i, err, "Remove activity command validation failed")
		return
	}

	name, err := getStringOption(i, 0)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get name option")
		return
	}

	err = data.ActivityRepo.RemoveActivity(name)
	if err != nil {
		handleCommandError(s, i, err, "Failed to remove activity")
		return
	}

	description := fmt.Sprintf("Successfully removed activity **%s**", name)
	embed := createSuccessEmbed("✅ Activity Removed Successfully", description, i)
	sendEmbedResponse(s, i, embed)
}
