package service

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Skill struct {
	Name       string
	Rank       int
	Level      int
	Experience int
}

type Activity struct {
	Name   string
	Rank   int
	Amount int
}

func CheckIfPlayerExists(username string) bool {
	url := fmt.Sprintf("https://secure.runescape.com/m=hiscore_oldschool/index_lite.ws?player=%s", username)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}

	return true
}

func FetchHiscore(username string) (map[string]Skill, map[string]Activity, error) {
	url := fmt.Sprintf("https://secure.runescape.com/m=hiscore_oldschool/index_lite.ws?player=%s", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch hiscore data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	lines := strings.Split(string(body), "\n")

	skillsOrder := []string{
		"Overall", "Attack", "Defence", "Strength", "Hitpoints", "Ranged", "Prayer", "Magic",
		"Cooking", "Woodcutting", "Fletching", "Fishing", "Firemaking", "Crafting", "Smithing",
		"Mining", "Herblore", "Agility", "Thieving", "Slayer", "Farming", "Runecrafting", "Hunter",
		"Construction",
	}

	activitiesOrder := []string{
		"League Points", "Deadman Points", "Bounty Hunter - Hunter", "Bounty Hunter - Rogue",
		"Bounty Hunter (Legacy) - Hunter", "Bounty Hunter (Legacy) - Rogue", "Clue Scrolls (all)",
		"Clue Scrolls (beginner)", "Clue Scrolls (easy)", "Clue Scrolls (medium)", "Clue Scrolls (hard)",
		"Clue Scrolls (elite)", "Clue Scrolls (master)", "LMS - Rank", "PvP Arena - Rank",
		"Soul Wars Zeal", "Rifts closed", "Colosseum Glory", "Abyssal Sire", "Alchemical Hydra",
		"Amoxliatl", "Araxxor", "Artio", "Barrows Chests", "Bryophyta", "Callisto", "Cal'varion",
		"Cerberus", "Chambers of Xeric", "Chambers of Xeric: Challenge Mode", "Chaos Elemental",
		"Chaos Fanatic", "Commander Zilyana", "Corporeal Beast", "Crazy Archaeologist",
		"Dagannoth Prime", "Dagannoth Rex", "Dagannoth Supreme", "Deranged Archaeologist",
		"Duke Sucellus", "General Graardor", "Giant Mole", "Grotesque Guardians", "Hespori",
		"Kalphite Queen", "King Black Dragon", "Kraken", "Kree'Arra", "K'ril Tsutsaroth",
		"Lunar Chests", "Mimic", "Nex", "Nightmare", "Phosani's Nightmare", "Obor",
		"Phantom Muspah", "Sarachnis", "Scorpia", "Scurrius", "Skotizo", "Sol Heredit", "Spindel",
		"Tempoross", "The Gauntlet", "The Corrupted Gauntlet", "The Hueycoatl", "The Leviathan",
		"The Whisperer", "Theatre of Blood", "Theatre of Blood: Hard Mode", "Thermonuclear Smoke Devil",
		"Tombs of Amascut", "Tombs of Amascut: Expert Mode", "TzKal-Zuk", "TzTok-Jad", "Vardorvis",
		"Venenatis", "Vet'ion", "Vorkath", "Wintertodt", "Zalcano", "Zulrah",
	}

	skills := make(map[string]Skill)
	activities := make(map[string]Activity)

	for i, line := range lines {
		if line == "" {
			continue
		}

		values := strings.Split(line, ",")
		if len(values) < 2 {
			continue
		}

		if i < len(skillsOrder) {
			skills[skillsOrder[i]] = Skill{
				Name:       skillsOrder[i],
				Rank:       atoi(values[0]),
				Level:      atoi(values[1]),
				Experience: atoi(values[2]),
			}
		} else if i-len(skillsOrder) < len(activitiesOrder) {
			activityIndex := i - len(skillsOrder)
			activities[activitiesOrder[activityIndex]] = Activity{
				Name:   activitiesOrder[activityIndex],
				Rank:   atoi(values[0]),
				Amount: atoi(values[1]),
			}
		}
	}

	return skills, activities, nil
}

func atoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}
