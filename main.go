package main

import (
	"github.com/abgeo/goclockify/src/context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	ui "github.com/gizak/termui/v3"
)

var (
	updateInterval = time.Second
)

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

func eventLoop(appContext *context.AppContext) {
	drawTicker := time.NewTicker(updateInterval).C

	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)

	uiEvents := ui.PollEvents()
	previousKey := ""

	for {
		select {
		case <-sigTerm:
			return
		case <-drawTicker:
			ui.Render(appContext.Grid)
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				appContext.Grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(appContext.Grid)

			// Workplaces events.
			case "a":
				appContext.View.Workplaces.ScrollUp()
				updateTimeEntries(appContext)
				ui.Render(appContext.View.Workplaces)
			case "z":
				appContext.View.Workplaces.ScrollDown()
				updateTimeEntries(appContext)
				ui.Render(appContext.View.Workplaces)

			// TimeEntries events.
			case "<MouseLeft>":
				payload := e.Payload.(ui.Mouse)
				appContext.View.TimeEntries.HandleClick(payload.X, payload.Y)
				ui.Render(appContext.View.TimeEntries)
			case "k", "<Up>", "<MouseWheelUp>":
				appContext.View.TimeEntries.ScrollUp()
				ui.Render(appContext.View.TimeEntries)
			case "j", "<Down>", "<MouseWheelDown>":
				appContext.View.TimeEntries.ScrollDown()
				ui.Render(appContext.View.TimeEntries)
			case "g", "<Home>":
				appContext.View.TimeEntries.ScrollTop()
				ui.Render(appContext.View.TimeEntries)
			case "G", "<End>":
				appContext.View.TimeEntries.ScrollBottom()
				ui.Render(appContext.View.TimeEntries)
			}

			if previousKey == e.ID {
				previousKey = ""
			} else {
				previousKey = e.ID
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
