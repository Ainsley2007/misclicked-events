package commands

import (
	"fmt"
	"misclicked-events/data"
	"misclicked-events/utils"

	"github.com/bwmarrin/discordgo"
)

func TrackedAccountsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	accounts, err := data.TrackedAccounts(i.GuildID, i.Member.User.ID)
	if err != nil {
		utils.RespondWithError(s, i, err)
		return
	}

	if len(accounts) == 0 {
		utils.RespondWithPrivateMessage(s, i, "You have no tracked accounts at the moment. Use `/track` to start tracking one!")
		return
	}

	currentCompetition := data.GetCurrentBoss(i.GuildID)
	response := "**🔍 Currently Tracked Accounts (Event: " + currentCompetition + "):**\n\n"
	if len(currentCompetition) == 0 {
		response = "**🔍 Currently Tracked Accounts:**\n\n"
		response += "*No event is currently running.*\n"
	}

	for _, account := range accounts {
		if len(currentCompetition) > 0 {
			activity, ok := account.Activities[currentCompetition]
			if ok {
				response += fmt.Sprintf(
					"🔹 **%s**\n   └ **%s KC (Current Event)**: `%d`\n",
					account.Name,
					activity.Name,
					activity.CurrentAmount-activity.StartAmount,
				)
			} else {
				response += fmt.Sprintf("🔹 **%s**\n   └ *Not participating in the current event*\n", account.Name)
			}
		} else {
			response += fmt.Sprintf("🔹 **%s**\n   └ *No active competition*\n", account.Name)
		}
	}

	utils.RespondWithMessage(s, i, "%s", response)
}
