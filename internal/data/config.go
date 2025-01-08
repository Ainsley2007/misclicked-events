package data

import (
	"encoding/json"
	"fmt"
	"os"
)

type BotConfig struct {
	CategoryChannelID string `json:"categoryChannelId"`
	HiscoreChannelID  string `json:"hiscoreChannelId"`
	HiscoreMessageID  string `json:"hiscoreMessageId"`
	RankingChannelID  string `json:"rankingChannelId"`
	RankingMessageID  string `json:"rankingMessageId"`
}

const configPath = "./assets/%s_config.json"

func SaveBotConfig(guildID string, botConfig BotConfig) error {
	data, err := json.MarshalIndent(botConfig, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal persons: %w", err)
	}

	file, err := os.Create(fmt.Sprintf(configPath, guildID))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func GetBotConfig(guildID string) (*BotConfig, error) {
	file, err := os.Open(fmt.Sprintf(configPath, guildID))
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	if stat.Size() == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config BotConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &config, nil
}

func UpdateChannelIDs(guildID string, newChannels BotConfig) error {
	config, err := GetBotConfig(guildID)
	if err != nil {
		return fmt.Errorf("failed to get bot config: %w", err)
	}

	// Update the channel IDs and clear message IDs
	config.CategoryChannelID = newChannels.CategoryChannelID
	config.HiscoreChannelID = newChannels.HiscoreChannelID
	config.RankingChannelID = newChannels.RankingChannelID
	config.HiscoreMessageID = ""
	config.RankingMessageID = ""

	return SaveBotConfig(guildID, *config)
}

func UpdateHiscoreMessageID(guildID, hiscoreMessageID string) error {
	config, err := GetBotConfig(guildID)
	if err != nil {
		return fmt.Errorf("failed to get bot config: %w", err)
	}

	config.HiscoreMessageID = hiscoreMessageID

	return SaveBotConfig(guildID, *config)
}

func UpdateRankingMessageID(guildID, rankingMessageID string) error {
	config, err := GetBotConfig(guildID)
	if err != nil {
		return fmt.Errorf("failed to get bot config: %w", err)
	}

	config.RankingMessageID = rankingMessageID

	return SaveBotConfig(guildID, *config)
}
