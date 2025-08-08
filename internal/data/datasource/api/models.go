package api

type HiscoreDataModel struct {
	Skills     []SkillModel    `json:"skills"`
	Activities []ActivityModel `json:"activities"`
}

type SkillModel struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Level int    `json:"level"`
	XP    int    `json:"xp"`
}

type ActivityModel struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Score int    `json:"score"`
}
