package constants

var Activities = map[string]ActivityDetails{
	"COLO": ActivityDetails{
		Threshold:     5,
		BossNames:     []string{"Sol Heredit"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_solheredit.png?2",
	},
	"Corp": ActivityDetails{
		Threshold:     25,
		BossNames:     []string{"Corporeal Beast"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_corporealbeast.png?2",
	},
	"Wildy": ActivityDetails{
		Threshold: 25,
		BossNames: []string{"Artio", "Callisto", "Cal'varion", "Vet'ion", "Venenatis", "Spindel"},
	},
	"COX": ActivityDetails{
		Threshold:     5,
		BossNames:     []string{"Chambers of Xeric", "Chambers of Xeric: Challenge Mode"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_chambersofxeric.png?2",
	},
	"Huey": ActivityDetails{
		Threshold:     25,
		BossNames:     []string{"The Hueycoatl"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_thehueycoatl.png?2",
	},
	"Inferno": ActivityDetails{
		Threshold:     5,
		BossNames:     []string{"TzKal-Zuk"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_tzkalzuk.png?2",
	},
	"Nex": ActivityDetails{
		Threshold:     25,
		BossNames:     []string{"Nex"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_nex.png?2",
	},
	"NM": ActivityDetails{
		Threshold:     25,
		BossNames:     []string{"Nightmare", "Phosani's Nightmare"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_nightmare.png?2",
	},
	"Sarachnis": ActivityDetails{
		Threshold:     25,
		BossNames:     []string{"Sarachnis"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_sarachnis.png?2",
	},
	"TOA": ActivityDetails{
		Threshold:     5,
		BossNames:     []string{"Tombs of Amascut", "Tombs of Amascut: Expert Mode"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_tombsofamascutexpertmode.png?2",
	},
	"TOB": ActivityDetails{
		Threshold:     5,
		BossNames:     []string{"Theatre of Blood", "Theatre of Blood: Hard Mode"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_theatreofblood.png?2",
	},
	"Zulrah": ActivityDetails{
		Threshold:     25,
		BossNames:     []string{"Zulrah"},
		BossThumbnail: "https://www.runescape.com/img/rsp777/game_icon_zulrah.png?2",
	},
}

type ActivityDetails struct {
	Threshold     int
	BossNames     []string
	BossThumbnail string
}
