package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

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
		epanic(err, "Where the csv")
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

func parseTeamsFromMatches() []Team {
	matchesFile, err := os.OpenFile("twoMatches.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		epanic(err, "Where the matches")
	}
	defer matchesFile.Close()

	teams := []Team{}

	scanner := bufio.NewScanner(matchesFile)
	for scanner.Scan() {
		name := scanner.Text()

		teams = append(teams, Team{
			ID:   teamToID(name),
			Name: name,
		})
	}

	return teams
}
