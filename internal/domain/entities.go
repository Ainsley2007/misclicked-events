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
	CurrentKingParticipant int64
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

// Account represents a participant with their OSRS accounts and participation
type Account struct {
	ID                int64
	DiscordID         string
	Points            int
	OSRSAccounts      []OSRSAccount
	BotmParticipation *BotmParticipation
	KotsParticipation *KotsParticipation
	BotmEnabled       bool
	KotsEnabled       bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// OSRSAccount represents an OSRS account linked to a participant
type OSRSAccount struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BotmParticipation represents participation in a BOTM competition
type BotmParticipation struct {
	ID            int64
	BotmID        int64
	StartAmount   int
	CurrentAmount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// KotsParticipation represents participation in a KOTS competition
type KotsParticipation struct {
	ID            int64
	KotsID        int64
	StartAmount   int
	CurrentAmount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
