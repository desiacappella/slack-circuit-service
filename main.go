package main

import (
	"sort"
)

func main() {
	teams := parseTeamsFromMirchi()

	for _, t := range teams {
		slackSendToChannel("team-"+t.ID, mirchiMsg(t.MirchiLink))
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
