package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// GetToken retrieves the bot token from an environment variable
func GetToken() string {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	// Get the bot token
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		fmt.Println("No bot token found! Set the DISCORD_BOT_TOKEN environment variable.")
		os.Exit(1)
	}

	return token
}
