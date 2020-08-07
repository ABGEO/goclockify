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
	"github.com/abgeo/goclockify/configs"
	w "github.com/abgeo/goclockify/internal/widgets"
	"io/ioutil"
	"net/http"
	"time"
)

// ClockifyService is a service to work with the Clockify API
type ClockifyService struct {
	BaseURL     string
	Config      *configs.Config
	Client      http.Client
	CurrentUser w.User
}

// NewClockifyService creates new Clockify service
func NewClockifyService(cnfg *configs.Config) (*ClockifyService, error) {
	service := &ClockifyService{
		BaseURL: "https://api.clockify.me/api/v1",
		Config:  cnfg,
		Client: http.Client{
			Timeout: time.Second * 5,
		},
	}

	currentUser, err := service.getCurrentUser()
	if err != nil || currentUser.ID == "" {
		return nil, fmt.Errorf("not able to authorize client, check your connection and if your Clockify API "+
			"token is set correctly.\nConfig file: %s", configs.FilePath)
	}

	service.CurrentUser = currentUser

	return service, nil
}

func (s *ClockifyService) getCurrentUser() (w.User, error) {
	body, err := s.get(s.BaseURL + "/user")
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

// GetWorkplaces gets all workspaces from the API
func (s *ClockifyService) GetWorkplaces() ([]w.Workplace, error) {
	body, err := s.get(s.BaseURL + "/workspaces")
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

// GetTimeEntries gets latest time entries from given workspace
func (s *ClockifyService) GetTimeEntries(workspaceID string) ([]w.TimeEntry, error) {
	body, err := s.get(
		fmt.Sprintf(
			"%s/workspaces/%s/user/%s/time-entries?hydrated=true&page-size=200",
			s.BaseURL,
			workspaceID,
			s.CurrentUser.ID,
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

// RemoveTimeEntry deletes given time entry
func (s *ClockifyService) DeleteTimeEntry(workspaceID string, id string) error {
	url := fmt.Sprintf(
		"%s/workspaces/%s/time-entries/%s",
		s.BaseURL,
		workspaceID,
		id,
	)

	res, err := s.doRequest(http.MethodDelete, url)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return errors.New("unable to delete selected time entry")
	}

	return nil
}

func (s *ClockifyService) get(url string) ([]byte, error) {
	res, err := s.doRequest(http.MethodGet, url)
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

func (s *ClockifyService) doRequest(method string, url string) (*http.Response, error) {
	client := s.Client

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", s.Config.ClockifyAPIToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
