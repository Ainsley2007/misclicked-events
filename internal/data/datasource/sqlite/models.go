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
	CurrentKingParticipant int64
	Streak                 int
	StartDate              time.Time
	EndDate                *time.Time
	Status                 string
}

type ParticipantModel struct {
	ID         int64
	ServerID   string
	DiscordID  string
	BotmPoints int
	KotsPoints int
}

type AccountModel struct {
	ID               int64
	ParticipantID    int64
	Username         string
	FailedFetchCount int
}

type BotmParticipationModel struct {
	ParticipantID int64
	BotmID        int64
	StartAmount   int
	CurrentAmount int
}

type KotsParticipationModel struct {
	ParticipantID int64
	KotsID        int64
	StartAmount   int
	CurrentAmount int
}
