package sqlite

import "time"

type ServerModel struct {
	ID   string
	Name string
}

type ConfigModel struct {
	ServerID          string
	RankingChannelID  string
	HiscoreChannelID  string
	CategoryChannelID string
	RankingMessageID  string
	HiscoreMessageID  string
}

type BotmModel struct {
	ID          int64
	ServerID    string
	CurrentBoss string
	Password    string
	Status      string
}

type KotsModel struct {
	ID                     int64
	ServerID               string
	CurrentSkill           string
	CurrentKingParticipant string
	Streak                 int
	StartDate              time.Time
	EndDate                *time.Time
	Status                 string
}

type ParticipantModel struct {
	DiscordID  string
	ServerID   string
	BotmPoints int
	KotsPoints int
}

type AccountModel struct {
	ID               int64
	ParticipantID    string
	Username         string
	FailedFetchCount int
}

type BotmParticipationModel struct {
	ParticipantID string
	BotmID        int64
	StartAmount   int
	CurrentAmount int
}

type KotsParticipationModel struct {
	ParticipantID string
	KotsID        int64
	StartAmount   int
	CurrentAmount int
}
