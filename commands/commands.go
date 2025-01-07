package commands

import (
	"fmt"
	"misclicked-events/utils"

	"github.com/bwmarrin/discordgo"
)

func RegisterCommands(s *discordgo.Session) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "setup-channels",
		Description: "setup channels to show competition results",
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
				Description: "Select a channels category",
				Required:    false,
			},
		},
	})
	if err != nil {
		utils.LogError("Error creating config command:", err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "track",
		Description: "Tracks killcount for this account and adds it to your discord user.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "username",
				Description: "The username to start tracking",
				Required:    true,
			},
		},
	})
	if err != nil {
		utils.LogError("Error creating track command:", err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "untrack",
		Description: "untracks an account",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "username",
				Description: "The username to stop tracking",
				Required:    true,
			},
		},
	})
	if err != nil {
		utils.LogError("Error creating untrack command:", err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "tracking",
		Description: "accounts you're currently tracking",
	})
	if err != nil {
		utils.LogError("Error creating tracking command:", err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "start",
		Description: "Select an activity to start",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "choice",
				Description: "Choose an activity",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "Colosseum", Value: "COLO"},
					{Name: "Corporeal beast", Value: "Corp"},
					{Name: "Wildy boss trio (Vet'ion, Callisto, Venenatis)", Value: "Wildy"},
					{Name: "COX", Value: "COX"},
					{Name: "Huey", Value: "Huey"},
					{Name: "Inferno", Value: "Inferno"},
					{Name: "Nex", Value: "Nex"},
					{Name: "NM + PNM", Value: "NM"},
					{Name: "Sarachnis", Value: "Sarachnis"},
					{Name: "TOA", Value: "TOA"},
					{Name: "TOB", Value: "TOB"},
					{Name: "Zulrah", Value: "Zulrah"},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "password",
				Description: "Set an activity password",
				Required:    true,
			},
		},
	})
	if err != nil {
		utils.LogError("Error creating start command:", err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
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
	})
	if err != nil {
		utils.LogError("Error creating end command:", err)
	} else {
		fmt.Println("Commands registered successfully.")
	}
}
