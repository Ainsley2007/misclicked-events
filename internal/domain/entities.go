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
	ID         int64
	ServerID   string
	ActivityID int64
	Password   string
	Status     string
}

type BotmWithActivity struct {
	ID         int64
	ServerID   string
	ActivityID int64
	Password   string
	Status     string
	Activity   *ActivityEntity
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
	ID           int64
	Name         string
	Type         string
	HiscoreNames []string
	Threshold    int
}

type Account struct {
	DiscordID         string
	BotmPoints        int
	KotsPoints        int
	OSRSAccounts      []OSRSAccount
	BotmParticipation *BotmParticipation
	KotsParticipation *KotsParticipation
}

type ParticipantWithAccountKC struct {
	DiscordID string
	Accounts  []AccountWithKC
}

type AccountWithKC struct {
	ID       int64
	Username string
	KCGained int
}

type OSRSAccount struct {
	ID   int64
	Name string
}

type BotmParticipation struct {
	AccountID     int64
	BotmID        int64
	StartAmount   int
	CurrentAmount int
}

type KotsParticipation struct {
	AccountID     int64
	KotsID        int64
	StartAmount   int
	CurrentAmount int
}
