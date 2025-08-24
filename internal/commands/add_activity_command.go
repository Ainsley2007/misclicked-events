package commands

import (
	"fmt"
	"misclicked-events/internal/data"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var AddActivityCommand = &discordgo.ApplicationCommand{
	Name:        "add-activity",
	Description: "Add a new activity to the system",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "The name of the activity",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "type",
			Description: "The type of activity",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Boss",
					Value: "boss",
				},
				{
					Name:  "Skill",
					Value: "skill",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "hiscore_names",
			Description: "Comma-separated list of hiscore names (optional - will use activity name if not provided)",
			Required:    false,
		},
	},
}

func HandleAddActivityCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := deferResponse(s, i, "add-activity"); err != nil {
		return
	}

	if err := validateRequiredOptions(i, 2, "add-activity"); err != nil {
		handleCommandError(s, i, err, "Add activity command validation failed")
		return
	}

	name, err := getStringOption(i, 0)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get name option")
		return
	}

	activityType, err := getStringOption(i, 1)
	if err != nil {
		handleCommandError(s, i, err, "Failed to get type option")
		return
	}

	var hiscoreNames []string
	hiscoreNamesStr, err := getStringOption(i, 2)
	if err == nil && hiscoreNamesStr != "" {
		hiscoreNames = strings.Split(hiscoreNamesStr, ",")
		for i, name := range hiscoreNames {
			hiscoreNames[i] = strings.TrimSpace(name)
		}
	}

	err = data.AddActivityUseCase.Execute(name, activityType, hiscoreNames)
	if err != nil {
		handleCommandError(s, i, err, "Failed to add activity")
		return
	}

	var hiscoreNamesText string
	if len(hiscoreNames) > 0 {
		hiscoreNamesText = fmt.Sprintf("\nHiscore names: %s", strings.Join(hiscoreNames, ", "))
	} else {
		hiscoreNamesText = "\nHiscore name: " + name
	}

	description := fmt.Sprintf("Successfully added activity **%s** (Type: %s)%s", name, activityType, hiscoreNamesText)
	embed := createSuccessEmbed("✅ Activity Added Successfully", description, i)
	sendEmbedResponse(s, i, embed)
}
