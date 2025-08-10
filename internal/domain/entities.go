package domain

import "time"

type Server struct {
	ID   string
	Name string
}

type Config struct {
	RankingChannelID  string
	HiscoreChannelID  string
	CategoryChannelID string
	RankingMessageID  string
	HiscoreMessageID  string
}

type Botm struct {
	ID          int64
	CurrentBoss string
	Password    string
	Status      string
}

type Kots struct {
	ID                     int64
	CurrentSkill           string
	CurrentKingParticipant string
	Streak                 int
	StartDate              time.Time
	EndDate                *time.Time
	Status                 string
}

type HiscoreData struct {
	Skills     []Skill
	Activities []Activity
}

type Skill struct {
	ID    int
	Name  string
	Rank  int
	Level int
	XP    int
}

type Activity struct {
	ID    int
	Name  string
	Rank  int
	Score int
}

type ActivityEntity struct {
	Name         string
	Type         string
	HiscoreNames []string
}

type Account struct {
	DiscordID         string
	BotmPoints        int
	KotsPoints        int
	OSRSAccounts      []OSRSAccount
	BotmParticipation *BotmParticipation
	KotsParticipation *KotsParticipation
}

type OSRSAccount struct {
	ID   int64
	Name string
}

type BotmParticipation struct {
	ParticipantID string
	BotmID        int64
	StartAmount   int
	CurrentAmount int
}

type KotsParticipation struct {
	ParticipantID string
	KotsID        int64
	StartAmount   int
	CurrentAmount int
}
