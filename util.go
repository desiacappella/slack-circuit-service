package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

func parseTeamsFromCaptains() []Team {
	teamsFile, err := os.OpenFile("captains.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer teamsFile.Close()

	r := csv.NewReader(teamsFile)

	teams := []Team{}

	// Ignore headers
	r.Read()
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		team := Team{
			Name:       record[0],
			ID:         teamToID(record[0]),
			University: record[1],
			Liaison: Person{
				record[2],
				record[3],
			},
			Captain: Person{
				Name:  record[4],
				Email: record[5],
			},
		}

		// Parse out officers
		for i := 6; i < len(record); i += 2 {
			if len(strings.TrimSpace(record[i])) == 0 {
				break
			}

			team.Officers = append(team.Officers, Person{Name: record[i], Email: record[i+1]})
		}

		teams = append(teams, team)
	}

	return teams
}

func parseTeamsFromMirchi() []Team {
	teamsFile, err := os.OpenFile("mirchi-teams.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer teamsFile.Close()

	r := csv.NewReader(teamsFile)

	teams := []Team{}

	// Ignore headers
	r.Read()
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		teams = append(teams, Team{
			Name:       record[0],
			ID:         teamToID(record[0]),
			MirchiLink: record[1],
		})
	}

	return teams
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
