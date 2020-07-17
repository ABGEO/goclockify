// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abgeo/goclockify/src/config"
	w "github.com/abgeo/goclockify/src/widgets"
	"io/ioutil"
	"net/http"
	"time"
)

type ClockifyService struct {
	BaseUrl     string
	Config      *config.Config
	Client      http.Client
	CurrentUser w.User
}

func NewClockifyService(cnfg *config.Config) (*ClockifyService, error) {
	service := &ClockifyService{
		BaseUrl: "https://api.clockify.me/api/v1/",
		Config:  cnfg,
		Client: http.Client{
			Timeout: time.Second * 5,
		},
	}

	currentUser, err := service.getCurrentUser()
	if err != nil || currentUser.ID == "" {
		return nil, errors.New(
			fmt.Sprintf("not able to authorize client, check your connection and if your Clockify API token is "+
				"set correctly.\nConfig file: %s", config.FilePath))
	}

	service.CurrentUser = currentUser

	return service, nil
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

func (self *ClockifyService) GetTimeEntries(workspaceId string) ([]w.TimeEntry, error) {
	body, err := self.get(
		fmt.Sprintf(
			"%s/workspaces/%s/user/%s/time-entries?hydrated=true&page-size=200",
			self.BaseUrl,
			workspaceId,
			self.CurrentUser.ID,
		),
	)
	if err != nil {
		return nil, err
	}

	var timeEntries []w.TimeEntry
	err = json.Unmarshal(body, &timeEntries)
	if err != nil {
		return nil, err
	}

	return timeEntries, nil
}

func (self *ClockifyService) get(url string) ([]byte, error) {
	res, err := self.doGet(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("not able to get resource, check your connection")
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

func (self *ClockifyService) doGet(url string) (*http.Response, error) {
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

	return res, nil
}

func (self *ClockifyService) getCurrentUser() (w.User, error) {
	body, err := self.get(self.BaseUrl + "/user")
	var user w.User
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}
