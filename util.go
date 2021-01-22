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
	Name       string
	University string
	Liaison    Person
	Captain    Person
	Officers   []Person
}

var unusualMapping = map[string]string{
	"New York Masti": "masti",
}

func parseTeams() []Team {
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
