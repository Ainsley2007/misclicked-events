package sqlite

import "time"

type ServerModel struct {
	ID   string
	Name string
}
type Config struct {
	ServerID          string
	RankingChannelID  string
	HiscoreChannelID  string
	CategoryChannelID string
	RankingMessageID  string
	HiscoreMessageID  string
}

type Botm struct {
	ID          int64
	ServerID    string
	CurrentBoss string
	Password    string
	Status      string
}

type Kots struct {
	ID                     int64
	ServerID               string
	CurrentSkill           string
	CurrentKingParticipant int64
	Streak                 int
	StartDate              time.Time
	EndDate                *time.Time
	Status                 string
}
