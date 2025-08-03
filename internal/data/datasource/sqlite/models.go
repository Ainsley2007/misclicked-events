package sqlite

type Config struct {
	ServerID          string
	RankingChannelID  string
	HiscoreChannelID  string
	CategoryChannelID string
	RankingMessageID  string
	HiscoreMessageID  string
}

type Competition struct {
	ServerID    string
	CurrentBoss string
	Password    string
}

type ServerModel struct {
	ID   string
	Name string
}
