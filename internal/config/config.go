package config

import (
	"os"

	"misclicked-events/internal/utils"

	"github.com/joho/godotenv"
)

func GetToken() string {
	err := godotenv.Load("./.env")
	if err != nil {
		utils.Error("Error loading .env file: %v", err)
		os.Exit(1)
	}

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		utils.Error("No bot token found! Set the DISCORD_BOT_TOKEN environment variable.")
		os.Exit(1)
	}

	return token
}
