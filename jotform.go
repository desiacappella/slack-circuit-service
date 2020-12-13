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

func jotformGetInactiveCaptains() {
	data := getSubmissions()

	fmt.Println("code", data["responseCode"].(float64))

	// Get people who have already graduated
	submissions := data["content"].([]interface{})

	fmt.Println("num responses", len(submissions))

	graduatedPeopleAnswers := make([]map[string]interface{}, 0)

	for _, x := range submissions {
		// Get graduation year
		resp := x.(map[string]interface{})

		answers := resp["answers"].(map[string]interface{})
		for _, a := range answers {
			ac := a.(map[string]interface{})
			if strings.Index(ac["name"].(string), "graduationYear") >= 0 {
				// found the right index
				year := ac["answer"].(string)
				fmt.Println("graduation year is", year)

				if y, _ := strconv.Atoi(year); y <= latestGradYear {
					// this foo be graduated
					graduatedPeopleAnswers = append(graduatedPeopleAnswers, answers)
				}
			}
		}
	}

	fmt.Println("we have", len(graduatedPeopleAnswers), "graduated people")
}
