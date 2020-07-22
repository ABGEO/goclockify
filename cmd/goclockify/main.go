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
	"github.com/docopt/docopt.go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	ui "github.com/gizak/termui/v3"
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

var (
	updateInterval      = time.Second
	showDashboard       = true
	showSingleTimeEntry = false
	showHelp            = false
)

func init() {
	_, err := docopt.ParseArgs(fmt.Sprintf(usage, configs.AppName, configs.Version), os.Args[1:], configs.Version)
	if err != nil {
		log.Fatalf("failed to parse arguments: %v", err)
	}
}

func conditionalRender(condition bool, element ui.Drawable) {
	if condition {
		ui.Render(element)
	}
}

func updateTimeEntries(appContext *context.AppContext) {
	selectedWorkspace, err := appContext.View.Workplaces.GetSelectedWorkplace()
	if err != nil {
		ui.Close()
		log.Fatalf("failed to select workplace: %v", err)
	}
	timeEntryItems, err := appContext.ClockifyService.GetTimeEntries(selectedWorkspace.ID)
	if err != nil {
		ui.Close()
		log.Fatalf("failed to get time entries: %v", err)
	}

	appContext.View.TimeEntries.UpdateData(timeEntryItems, selectedWorkspace)
}

func contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}

	return false
}

func eventLoop(appContext *context.AppContext) {
	drawTicker := time.NewTicker(updateInterval).C

	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)

	uiEvents := ui.PollEvents()

	for {
		select {
		case <-sigTerm:
			return
		case <-drawTicker:
			if showSingleTimeEntry {
				appContext.View.TimeEntry.UpdateTable()
				ui.Render(appContext.View.TimeEntry)
			}

			conditionalRender(showDashboard, appContext.Grid)
		case e := <-uiEvents:
			switch e.ID {
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				ui.Clear()

				if showDashboard {
					appContext.Grid.SetRect(0, 0, payload.Width, payload.Height)
					ui.Render(appContext.Grid)
				}

				if showSingleTimeEntry {
					appContext.View.TimeEntry.SetRect(0, 0, payload.Width, payload.Height)
					ui.Render(appContext.View.TimeEntry)
				}

				if showHelp {
					appContext.View.Help.SetRect(0, 0, payload.Width, payload.Height)
					ui.Render(appContext.View.Help)
				}

			case "<MouseLeft>":
				payload := e.Payload.(ui.Mouse)
				appContext.View.TimeEntries.HandleClick(payload.X, payload.Y)
				conditionalRender(showDashboard, appContext.View.TimeEntries)
			}

			// Configurable key mapping
			switch {
			case contains(e.ID, appContext.Config.KeyMapping.Other.Quit):
				return
			case contains(e.ID, appContext.Config.KeyMapping.Other.Help):
				showHelp = !showHelp
				showDashboard = !showHelp

				if showHelp {
					ui.Clear()
					terminalWidth, terminalHeight := ui.TerminalDimensions()
					appContext.View.Help.SetRect(0, 0, terminalWidth, terminalHeight)
					ui.Render(appContext.View.Help)
				}
			case contains(e.ID, appContext.Config.KeyMapping.Other.CloseWindow):
				if !showDashboard {
					showDashboard = true
					showSingleTimeEntry = false
					showHelp = false
					ui.Clear()
					ui.Render(appContext.Grid)
				}

			// Workplaces events.
			case contains(e.ID, appContext.Config.KeyMapping.Workspace.NavigationUp):
				appContext.View.Workplaces.ScrollUp()
				updateTimeEntries(appContext)
				conditionalRender(showDashboard, appContext.View.Workplaces)
			case contains(e.ID, appContext.Config.KeyMapping.Workspace.NavigationDown):
				appContext.View.Workplaces.ScrollDown()
				updateTimeEntries(appContext)
				conditionalRender(showDashboard, appContext.View.Workplaces)

			// TimeEntries events.
			case contains(e.ID, appContext.Config.KeyMapping.TimeEntries.NavigationUp):
				appContext.View.TimeEntries.ScrollUp()
				conditionalRender(showDashboard, appContext.View.TimeEntries)
			case contains(e.ID, appContext.Config.KeyMapping.TimeEntries.NavigationDown):
				appContext.View.TimeEntries.ScrollDown()
				conditionalRender(showDashboard, appContext.View.TimeEntries)
			case contains(e.ID, appContext.Config.KeyMapping.TimeEntries.NavigationToTop):
				appContext.View.TimeEntries.ScrollTop()
				conditionalRender(showDashboard, appContext.View.TimeEntries)
			case contains(e.ID, appContext.Config.KeyMapping.TimeEntries.NavigationToBottom):
				appContext.View.TimeEntries.ScrollBottom()
				conditionalRender(showDashboard, appContext.View.TimeEntries)
			case contains(e.ID, appContext.Config.KeyMapping.TimeEntries.NavigationSelect):
				showSingleTimeEntry = !showSingleTimeEntry
				showDashboard = !showSingleTimeEntry

				if timeEntry, err := appContext.View.TimeEntries.GetSelectedTimeEntry(); showSingleTimeEntry && err == nil {
					appContext.View.TimeEntry.SetTimeEntry(timeEntry)
					terminalWidth, terminalHeight := ui.TerminalDimensions()
					appContext.View.TimeEntry.SetRect(0, 0, terminalWidth, terminalHeight)

					ui.Clear()
					ui.Render(appContext.View.TimeEntry)
				}
			}
		}
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
		log.Fatalf("failed to create AppContext: %v", err)
	}

	eventLoop(appContext)
}
