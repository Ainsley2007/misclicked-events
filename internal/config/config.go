package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetToken() string {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		fmt.Println("No bot token found! Set the DISCORD_BOT_TOKEN environment variable.")
		os.Exit(1)
	}

	return token
}
