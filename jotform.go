package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

const jotformKey = "jotformApiKey"
const jotformCache = "jotform.json"

// /form/{formID}/submissions?apiKey={apiKey}"
const jotformURL = "https://api.jotform.com"

// TODO get this programmatically via /user/forms/
const formID = "92310716200139"

// The last year that we want to purge, inclusive
const latestGradYear = 2020

func getSubmissions() []jotformSubm {
	data, err := ioutil.ReadFile(jotformCache)
	if err != nil {
		// Try API call
		client := &http.Client{}

		v := url.Values{}
		v.Set("limit", "1000")

		req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s",
			strings.Join([]string{jotformURL, "form", formID, "submissions"}, "/"),
			v.Encode()), nil)
		if err != nil {
			epanic(err, "could not create request")
		}

		req.Header.Add("APIKEY", getToken(jotformKey))
		resp, err := client.Do(req)
		if err != nil {
			epanic(err, "could not get")
		}
		defer resp.Body.Close()

		data, err = ioutil.ReadAll(resp.Body)

		ioutil.WriteFile(jotformCache, data, 0644)
	}

	// Unmarshal
	processedBody := make(map[string]interface{})
	err = json.Unmarshal(data, &processedBody)

	// Sanitize and extract
	if processedBody["responseCode"].(float64) != 200 {
		epanic(fmt.Errorf("Failed to get submissions"), "Failed to get submissions")
	}

	submissions := processedBody["content"].([]interface{})

	var subms []jotformSubm

	for _, x := range submissions {
		resp := x.(map[string]interface{})

		answers := resp["answers"].(map[string]interface{})
		var subm jotformSubm

		for _, a := range answers {
			ac := a.(map[string]interface{})
			if strings.HasPrefix(ac["name"].(string), "name") {
				subm.Name = ac["prettyFormat"].(string)
			} else if strings.HasPrefix(ac["name"].(string), "yourEmail") {
				subm.Email = ac["answer"].(string)
			} else if strings.HasPrefix(ac["name"].(string), "groupName") {
				subm.TeamName = ac["answer"].(string)
			} else if strings.HasPrefix(ac["name"].(string), "graduationYear") {
				// found the right index
				year, err := strconv.Atoi(strings.TrimSpace(ac["answer"].(string)))
				if err == nil {
					subm.GradYear = year
				}
			}
		}

		subms = append(subms, subm)
	}

	return subms
}

type jotformSubm struct {
	Name     string
	Email    string
	GradYear int
	TeamName string
}

func jotformGetInactiveCaptains() []jotformSubm {
	submissions := getSubmissions()

	var graduatedPeople []jotformSubm

	for _, s := range submissions {
		if s.GradYear <= latestGradYear {
			graduatedPeople = append(graduatedPeople, s)
		}
	}

	fmt.Println("we have", len(graduatedPeople), "graduated people")

	return graduatedPeople
}

func jotformGetTeams() []Team {
	submissions := getSubmissions()

	// Go through and create a map of team ID's to players. Then convert that into teams
	// id -> name -> email
	membership := make(map[string]map[string]string)

	for _, s := range submissions {
		if len(s.TeamName) <= 0 || s.GradYear <= latestGradYear {
			continue
		}

		id := teamToID(s.TeamName)

		if _, ok := membership[id]; !ok {
			membership[id] = make(map[string]string)
		}

		// Prefer .edu emails
		email, ok := membership[id][s.Name]
		if !ok || !strings.HasSuffix(email, ".edu") {
			membership[id][s.Name] = s.Email
		}
	}

	var teams []Team

	for teamName, playerSet := range membership {
		team := Team{Name: teamName, Officers: make([]Person, len(playerSet))}

		i := 0
		for name, email := range playerSet {
			team.Officers[i] = Person{Name: name, Email: email}
			i++
		}

		teams = append(teams, team)
	}

	sort.SliceStable(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})
	return teams
}

// For each team add officers
func jotformUpdateOfficers(teams *[]Team) {
	submissions := getSubmissions()

	for _, s := range submissions {
		if s.GradYear <= latestGradYear {
			continue
		}

		foundTeam := getTeamFromName(s.TeamName, *teams)
		if foundTeam < 0 {
			fmt.Println("Did not map to a team:", s.TeamName)
			continue
		}

		(*teams)[foundTeam].Officers = append((*teams)[foundTeam].Officers, Person{
			Name:  s.Name,
			Email: s.Email,
		})
	}
}
