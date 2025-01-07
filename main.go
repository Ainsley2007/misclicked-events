package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"misclicked-events/commands"
	"misclicked-events/config"
	"misclicked-events/handlers"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := config.GetToken()
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(handlers.InteractionCreateHandler)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	commands.RegisterCommands(dg)

	handlers.UpdateBOTMHiscores(dg)

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	dg.Close()
}
