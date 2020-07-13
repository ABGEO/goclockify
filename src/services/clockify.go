package services

import (
	"encoding/json"
	"errors"
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
	CurrentUser User
}

type User struct {
	ID               string
	Email            string
	Name             string
	ProfilePicture   string
	ActiveWorkspace  string
	DefaultWorkspace string
}

func NewClockifyService(config *config.Config) (*ClockifyService, error) {
	service := &ClockifyService{
		BaseUrl: "https://api.clockify.me/api/v1/",
		Config:  config,
		Client: http.Client{
			Timeout: time.Second * 5,
		},
	}

	currentUser, err := service.getCurrentUser()
	if err != nil || currentUser.ID == "" {
		return nil, errors.New("not able to authorize client, check your connection and if your Clockify API token is set correctly")
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

func (self *ClockifyService) getCurrentUser() (User, error) {
	body, err := self.get(self.BaseUrl + "/user")
	var user User
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}
