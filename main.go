package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/desiacappella/slack-circuit-service/slackservice"
)

func collect(s *slackservice.SlackService) {
	channels := s.GetChannelsByPrefix(slackservice.PrefixCircuit)
	jsonString, err := json.MarshalIndent(channels, "", "\t")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile("circuit-ids.json", jsonString, 0644)
	if err != nil {
		log.Fatalf("Failed to write to JSON file: %v", err)
	}
}

func archive(s *slackservice.SlackService) {
	jsonString, err := os.ReadFile("circuit-ids.json")
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	channels := make([]string, 0)
	err = json.Unmarshal(jsonString, &channels)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON string: %v", err)
	}

	s.ArchiveAllChannels(channels)
}

func main() {
	s := slackservice.NewSlackService()

	if os.Args[1] == "collect" {
		collect(s)
	} else if os.Args[1] == "archive" {
		archive(s)
	}
}

/*
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
*/
