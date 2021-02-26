package main

import (
	"fmt"
	"sort"
)

func main() {
	teams := parseTeamsFromMatches()
	filtered := filterTeams(teams, []string{"dhun", "dhunki", "basmati-beats", "hum", "anokha"})

	fmt.Println(len(filtered))

	for _, t := range filtered {
		slackSendToChannel("team-"+t.ID, typePrivateChannel, extension2Msg)
		// slackSendToChannel("asa-bot-test-channel", typePublicChannel, extensionMsg)
	}
}

func createNewChannels() {
	teams := parseTeamsFromMirchi()
	jotformUpdateOfficers(&teams)

	sort.SliceStable(teams, func(i int, j int) bool {
		return teams[i].ID < teams[j].ID
	})

	compTeams := parseTeamsFromCaptains()

	filteredTeams := make([]Team, len(teams)-len(compTeams))
	i := 0
	for _, t := range teams {
		add := true
		for _, ct := range compTeams {
			if ct.ID == t.ID {
				add = false
				break
			}
		}
		if add {
			filteredTeams[i] = t
			i++
		}
	}

	slackTeamChannels(filteredTeams, false)
}
