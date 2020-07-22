// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package configs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OpenPeeDeeP/xdg"
	"io/ioutil"
	"os"
	"path/filepath"
	fp "path/filepath"
)

// AppName Application Name
// Version Application Version
const (
	AppName = "goclockify"
	Version = "1.0.0"
)

// FilePath Configuration file path
var (
	FilePath = xdg.New("abgeo", AppName).QueryConfig("config")
)

// WorkspaceKeyMapping is a structure for workspaces key mapping
type WorkspaceKeyMapping struct {
	NavigationUp   []string `json:"nav_up"`
	NavigationDown []string `json:"nav_down"`
}

// TimeEntriesKeyMapping is a structure for general key mapping
type TimeEntriesKeyMapping struct {
	NavigationUp       []string `json:"nav_up"`
	NavigationDown     []string `json:"nav_down"`
	NavigationToTop    []string `json:"nav_to_top"`
	NavigationToBottom []string `json:"nav_to_bottom"`
	NavigationSelect   []string `json:"nav_select"`
}

// OtherKeyMapping is a structure for general key mapping
type OtherKeyMapping struct {
	Quit        []string `json:"quit"`
	CloseWindow []string `json:"close_window"`
	Help        []string `json:"help"`
}

// KeyMapping the union of WorkspaceKeyMapping, TimeEntriesKeyMapping and OtherKeyMapping
type KeyMapping struct {
	Workspace   WorkspaceKeyMapping   `json:"workspace"`
	TimeEntries TimeEntriesKeyMapping `json:"time_entries"`
	Other       OtherKeyMapping       `json:"other"`
}

// Config structure type
type Config struct {
	ClockifyAPIToken string     `json:"clockify_api_token"`
	KeyMapping       KeyMapping `json:"key_mapping"`
}

// NewConfig creates the new config object from FilePath content
func NewConfig() (cfg *Config, err error) {
	if "" == FilePath {
		FilePath = filepath.Join(filepath.Join(os.Getenv("HOME"), ".config"), "abgeo/goclockify/config")
	}

	cfg = &Config{
		ClockifyAPIToken: "",
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

// CreateConfigFile creates the default config file
func CreateConfigFile() (file *os.File, err error) {
	if _, err := os.Stat(FilePath); os.IsNotExist(err) {
		err := os.MkdirAll(fp.Dir(FilePath), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	config := Config{
		ClockifyAPIToken: "",
		KeyMapping: KeyMapping{
			Workspace: WorkspaceKeyMapping{
				NavigationUp:   []string{"a"},
				NavigationDown: []string{"z"},
			},
			TimeEntries: TimeEntriesKeyMapping{
				NavigationUp:       []string{"k", "<Up>", "<MouseWheelUp>"},
				NavigationDown:     []string{"j", "<Down>", "<MouseWheelDown>"},
				NavigationToTop:    []string{"g", "<Home>"},
				NavigationToBottom: []string{"G", "<End>"},
				NavigationSelect:   []string{"<Enter>"},
			},
			Other: OtherKeyMapping{
				Quit:        []string{"q", "<C-c>"},
				CloseWindow: []string{"<Escape>"},
				Help:        []string{"<F1>", "?"},
			},
		},
	}
	payload, err := configToJSON(config)
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

func configToJSON(config Config) ([]byte, error) {
	b := &bytes.Buffer{}
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	err := e.Encode(config)

	var b2 bytes.Buffer
	json.Indent(&b2, b.Bytes(), "", "  ")

	return b2.Bytes(), err
}
