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
	ID          int64
	ServerID    string
	DiscordID   string
	Points      int
	BotmEnabled bool
	KotsEnabled bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AccountModel struct {
	ID            int64
	ParticipantID int64
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type BotmParticipationModel struct {
	ID            int64
	ParticipantID int64
	BotmID        int64
	StartAmount   int
	CurrentAmount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type KotsParticipationModel struct {
	ID            int64
	ParticipantID int64
	KotsID        int64
	StartAmount   int
	CurrentAmount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
