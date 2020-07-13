package config

import (
	"encoding/json"
	"fmt"
	"github.com/OpenPeeDeeP/xdg"
	"io/ioutil"
	"os"
	fp "path/filepath"
)

type Config struct {
	ClockifyApiToken string `json:"clockify_api_token"`
}

func NewConfig() (*Config, error) {
	filepath := xdg.New("abgeo", "goclockify").QueryConfig("config")
	cfg := Config{
		ClockifyApiToken: "",
	}

	file, err := os.Open(filepath)
	if err != nil {
		file, err = CreateConfigFile(filepath)
		if err != nil {
			return &cfg, fmt.Errorf("couldn't open the goclockify config file: (%v)", err)
		}
	}

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return &cfg, fmt.Errorf("the goclockify config file isn't valid json: (%v)", err)
	}

	return &cfg, nil
}

func CreateConfigFile(filepath string) (*os.File, error) {
	filepath = fmt.Sprintf("%s/abgeo/goclockify/config", xdg.ConfigHome())

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		err := os.MkdirAll(fp.Dir(filepath), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	payload := "{\"clockify_api_token\": \"\"}"
	err := ioutil.WriteFile(filepath, []byte(payload), 0755)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}