package services

import (
	"encoding/json"
	"github.com/abgeo/goclockify/src/config"
	w "github.com/abgeo/goclockify/src/widgets"
	"io/ioutil"
	"net/http"
	"time"
)

type ClockifyService struct {
	BaseUrl string
	Config  *config.Config
	Client  http.Client
}

func NewClockifyService(config *config.Config) (*ClockifyService, error) {
	return &ClockifyService{
		BaseUrl: "https://api.clockify.me/api/v1/",
		Config:  config,
		Client: http.Client{
			Timeout: time.Second * 5,
		},
	}, nil
}

func (self *ClockifyService) GetWorkplaces() ([]w.Workplace, error) {
	body, err := self.get(self.BaseUrl + "/workspaces")
	if err != nil {
		return nil, err
	}

	var workplaces []w.Workplace
	err = json.Unmarshal(body, &workplaces)
	if err != nil {
		return nil, err
	}

	return workplaces, nil
}

func (self *ClockifyService) get(url string) ([]byte, error) {
	spaceClient := self.Client

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", self.Config.ClockifyApiToken)

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
