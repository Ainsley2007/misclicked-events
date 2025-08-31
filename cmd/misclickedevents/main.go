package main

import (
	"os"
	"os/signal"
	"syscall"

	"misclicked-events/internal/commands"
	"misclicked-events/internal/config"
	"misclicked-events/internal/data"
	"misclicked-events/internal/handlers"
	"misclicked-events/internal/services"
	"misclicked-events/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {
	if err := data.Init("./data.db"); err != nil {
		utils.Error("could not init data layer: %v", err)
		return
	}

	token := config.GetToken()
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		utils.Error("error creating Discord session: %v", err)
		return
	}

	dg.AddHandler(handlers.InteractionCreateHandler)

	readyFn := handlers.MakeReadyHandler()
	dg.AddHandler(readyFn)

	err = dg.Open()
	if err != nil {
		utils.Error("error opening connection: %v", err)
		return
	}

	commands.RegisterCommands(dg, true)

	// Initialize and start the hiscore service
	hiscoreService := services.NewHiscoreService(dg, data.UpdateHiscoresUseCase)
	hiscoreService.Start()

	utils.Info("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Gracefully stop the hiscore service
	hiscoreService.Stop()
	dg.Close()
}
