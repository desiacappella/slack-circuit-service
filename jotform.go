package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

func getSubmissions() map[string]interface{} {
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

	return processedBody
}

type jotformSubm struct {
	Name     string
	Email    string
	GradYear int
}

func jotformGetInactiveCaptains() []jotformSubm {
	data := getSubmissions()

	if data["responseCode"].(float64) != 200 {
		epanic(fmt.Errorf("Failed to get submissions"), "Failed to get submissions")
	}

	// Get people who have already graduated
	submissions := data["content"].([]interface{})

	// graduatedPeopleAnswers := make([]map[string]interface{}, 0)
	var graduatedPeople []jotformSubm

	for _, x := range submissions {
		// Get graduation year
		resp := x.(map[string]interface{})

		answers := resp["answers"].(map[string]interface{})
		graduated := false
		var subm jotformSubm

		for _, a := range answers {
			ac := a.(map[string]interface{})
			if strings.Index(ac["name"].(string), "name") >= 0 {
				subm.Name = ac["prettyFormat"].(string)
			} else if strings.Index(ac["name"].(string), "yourEmail") >= 0 {
				subm.Email = ac["answer"].(string)
			} else if strings.Index(ac["name"].(string), "graduationYear") >= 0 {
				// found the right index
				year, err := strconv.Atoi(strings.TrimSpace(ac["answer"].(string)))
				if err == nil {
					subm.GradYear = year
					graduated = year <= latestGradYear
				}
			}
		}

		if graduated {
			graduatedPeople = append(graduatedPeople, subm)
		}
	}

	fmt.Println("we have", len(graduatedPeople), "graduated people")

	return graduatedPeople
}
