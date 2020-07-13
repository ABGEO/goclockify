package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	baseUrl = "https://api.clockify.me/api/v1/"
	apiKey  = "" // TODO: Move to config.
)

type Workplace struct {
	ID   string
	Name string
}

func GetWorkplaces() ([]Workplace, error) {
	body, err := get(baseUrl + "/workspaces")
	if err != nil {
		return nil, err
	}

	var workplaces []Workplace
	err = json.Unmarshal(body, &workplaces)
	if err != nil {
		return nil, err
	}

	return workplaces, nil
}

func get(url string) ([]byte, error) {
	spaceClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", apiKey)

	res, err := spaceClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
