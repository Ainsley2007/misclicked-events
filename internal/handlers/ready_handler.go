package handlers

import (
	"misclicked-events/internal/data"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func MakeReadyHandler() func(*discordgo.Session, *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		for _, g := range r.Guilds {
			name := g.Name
			if name == "" {
				full, err := s.Guild(g.ID)
				if err == nil {
					name = full.Name
				}
			}
			if err := data.ServerRepo.RegisterServer(g.ID, name); err != nil {
				utils.Error("could not register server %q (%s): %v", name, g.ID, err)
			}
		}
	}
}
