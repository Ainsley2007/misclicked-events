package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"misclicked-events/internal/commands"
	"misclicked-events/internal/config"
	"misclicked-events/internal/data"
	"misclicked-events/internal/handlers"

	"github.com/bwmarrin/discordgo"
)

func main() {
	if err := data.Init("./data.db"); err != nil {
		fmt.Println("could not init data layer:", err)
		return
	}

	token := config.GetToken()
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(handlers.InteractionCreateHandler)

	readyFn := handlers.MakeReadyHandler()
	dg.AddHandler(readyFn)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	commands.RegisterCommands(dg, false)

	commands.UpdateBOTMHiscores(dg)

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	dg.Close()
}
