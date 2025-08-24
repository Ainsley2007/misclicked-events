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
	ID         int64
	ServerID   string
	ActivityID int64
	Password   string
	Status     string
}

type BotmWithActivityModel struct {
	ID         int64
	ServerID   string
	ActivityID int64
	Password   string
	Status     string
	Activity   *ActivityModel
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
	AccountID     int64
	BotmID        int64
	StartAmount   int
	CurrentAmount int
}

type KotsParticipationModel struct {
	AccountID     int64
	KotsID        int64
	StartAmount   int
	CurrentAmount int
}

type ActivityModel struct {
	ID           int64
	Name         string
	Type         string
	HiscoreNames string
}
