package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func epanic(err error, msg string) {
	panic(fmt.Sprintf("%v, %s", err, msg))
}

// This token needs to be acquired after authenticating with the app, on the app's home page
func getToken(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Unable to find user token")
	}
	return strings.TrimSpace(string(b))
}

// Person is
type Person struct {
	Name  string
	Email string
}

// Team is
type Team struct {
	ID         string
	Name       string
	University string
	Liaison    Person
	Captain    Person
	Officers   []Person
	MirchiLink string
}

var unusualMapping = map[string]string{
	"New York Masti":          "masti",
	"Awaaz":                   "uw-awaaz",
	"Nuttin But V.O.C.A.L.S.": "nbv",
	"Madhura SJSU":            "madhura",
	"UC Davis Jhankaar":       "jhankaar",
}

func teamToID(team string) (id string) {
	id, ok := unusualMapping[team]
	if !ok {
		replaced := strings.ToLower(team)
		for _, s := range []string{"a cappella", "acappella", "acapella"} {
			replaced = strings.ReplaceAll(replaced, s, "")
		}
		return strings.ReplaceAll(strings.TrimSpace(strings.Split(replaced, ":")[0]), " ", "-")
	}
	return id
}

func getTeamFromName(name string, teams []Team) int {
	teamID := teamToID(name)

	foundTeam := -1
	for i, t := range teams {
		if t.ID == teamID {
			foundTeam = i
			break
		}
	}

	if foundTeam < 0 {
		// Try cutting off the first word, could be university name
		teamID = teamToID(name[strings.Index(name, " ")+1:])

		for i, t := range teams {
			if t.ID == teamID {
				foundTeam = i
				break
			}
		}
	}

	return foundTeam
}

func filterTeams(teams []Team, excludeIDs []string) []Team {
	filtered := []Team{}

	for _, t := range teams {
		exclude := false
		for _, n := range excludeIDs {
			if t.ID == n {
				exclude = true
				break
			}
		}
		if !exclude {
			filtered = append(filtered, t)
		}
	}

	return filtered
}
