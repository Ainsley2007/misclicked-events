package utils

import "github.com/bwmarrin/discordgo"

func IsAdmin(i *discordgo.InteractionCreate) bool {
	return i.Member.Permissions&discordgo.PermissionAdministrator != 0 ||
		i.Member.Permissions&discordgo.PermissionManageServer != 0
}
