package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// sentryProjects
type Organization struct {
	Id   string `json:"id"`
	Slug string `json:"slug"`
}

type Project struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Slug         string       `json:"slug"`
	Status       string       `json:"status"`
	Platform     string       `json:"platform"`
	Features     []string     `json:"features"`
	Organization Organization `json:"organization"`
}

// sentryIssues
type Stats struct {
	Stats24H [][]int `json:"24h"`
}

type Issue struct {
	Title     string `json:"title"`
	ShortId   string `json:"shortId"`
	PermaLink string `json:"permalink"`
	Stats     Stats  `json:"stats"`
}

var SENTRY_TOKEN string = getEnv("SENTRY_TOKEN", "UNDEFINED")

func sentryProjectList() (project []Project, err error) {
	client := http.Client{}
	req, _ := http.NewRequest("GET", "https://sentry.io/api/0/projects/", nil)

	req.Header = http.Header{
		"Authorization": []string{"Bearer " + SENTRY_TOKEN},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}

	// Read response and convert it to readable form
	r_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Response parse error\n", err)
	}

	// Parse response to predefined struct
	projects := []Project{}
	jsonErr := json.Unmarshal(r_body, &projects)

	// debug output
	// fmt.Printf("%#v", projects)
	return projects, jsonErr
}

func sentryProjectStats(organization string, project string) (stats [][]int, err error) {
	/*
		curl -H 'Authorization: Bearer TOKEN' \
			'https://sentry.io/api/0/projects/ORGANIZATION_SLUG/PROJECT_SLUG/stats/' | jq
	*/
	hour_ago := time.Now().Add(time.Duration(-1) * time.Hour)
	h_start := time.Date(hour_ago.Year(), hour_ago.Month(), hour_ago.Day(), hour_ago.Hour(), 0, 0, 0, time.UTC).Unix()
	h_end := time.Date(hour_ago.Year(), hour_ago.Month(), hour_ago.Day(), hour_ago.Hour(), 59, 0, 0, time.UTC).Unix()

	url := fmt.Sprintf("https://sentry.io/api/0/projects/%s/%s/stats/?stat=received&resolution=1h&since=%d&until=%d", organization, project, h_start, h_end)
	log.Printf("Parsing status from project %s/%s", organization, project)

	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header = http.Header{
		"Authorization": []string{"Bearer " + SENTRY_TOKEN},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}

	// Read response and convert it to readable form
	r_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Response parse error\n", err)
	}

	// Parse response to predefined struct
	jsonErr := json.Unmarshal(r_body, &stats)
	return stats, jsonErr

}

func sentryProjectIssues(organization string, project string) (issue []Issue, err error) {
	/*
		curl -H 'Authorization: Bearer TOKEN' \
			'https://sentry.io/api/0/projects/ORGANIZATION_SLUG/PROJECT_SLUG/issues/' | jq
	*/
	url := fmt.Sprintf("https://sentry.io/api/0/projects/%s/%s/issues/", organization, project)
	log.Printf("Parsing isuses from project %s/%s", organization, project)

	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header = http.Header{
		"Authorization": []string{"Bearer " + SENTRY_TOKEN},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}

	// Read response and convert it to readable form
	r_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Response parse error\n", err)
	}

	// Parse response to predefined struct
	issues := []Issue{}
	jsonErr := json.Unmarshal(r_body, &issues)

	// debug output
	// log.Println("Request status code: ", res.StatusCode)
	// fmt.Printf("%#v", issues[1])
	return issues, jsonErr

}
