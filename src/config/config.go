// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package config

import (
	"encoding/json"
	"fmt"
	"github.com/OpenPeeDeeP/xdg"
	"io/ioutil"
	"os"
	"path/filepath"
	fp "path/filepath"
)

const (
	AppName = "goclockify"
	Version = "1.0.0"
)

var (
	FilePath = xdg.New("abgeo", AppName).QueryConfig("config")
)

type Config struct {
	ClockifyApiToken string `json:"clockify_api_token"`
}

func NewConfig() (cfg *Config, err error) {
	if "" == FilePath {
		FilePath = filepath.Join(filepath.Join(os.Getenv("HOME"), ".config"), "abgeo/goclockify/config")
	}

	cfg = &Config{
		ClockifyApiToken: "",
	}

	file, err := os.Open(FilePath)
	if err != nil {
		file, err = CreateConfigFile()
		if err != nil {
			return cfg, fmt.Errorf("couldn't open the goclockify config file: (%v)", err)
		}
	}

	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("the goclockify config file isn't valid json: (%v)", err)
	}

	return cfg, nil
}

func CreateConfigFile() (file *os.File, err error) {
	if _, err := os.Stat(FilePath); os.IsNotExist(err) {
		err := os.MkdirAll(fp.Dir(FilePath), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	payload, err := json.Marshal(Config{})
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(FilePath, payload, 0755)
	if err != nil {
		return nil, err
	}

	file, err = os.Open(FilePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
