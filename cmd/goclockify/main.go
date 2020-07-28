// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package main

import (
	"fmt"
	"github.com/abgeo/goclockify/configs"
	"github.com/abgeo/goclockify/internal/context"
	"github.com/abgeo/goclockify/internal/handlers"
	"github.com/docopt/docopt.go"
	ui "github.com/gizak/termui/v3"
	"log"
	"os"
)

const (
	usage = `%[1]s - a terminal based client for Clockify time tracker

v%[2]s

Usage: %[1]s [options]

Options:
  -h, --help  Show this screen.
  --version   Print Application version.

`
)

func init() {
	_, err := docopt.ParseArgs(fmt.Sprintf(usage, configs.AppName, configs.Version), os.Args[1:], configs.Version)
	if err != nil {
		log.Fatalf("failed to parse arguments: %v", err)
	}
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize Termui: %v", err)
	}
	defer ui.Close()

	appContext, err := context.CreateAppContext()
	if err != nil {
		ui.Close()
		log.Fatal(err)
	}

	handlers.Initialize(appContext)
}
