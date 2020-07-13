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
			}

			switch e.ID {
			case "?":
				ui.Render(appContext.Grid)
			case "<Resize>":
				ui.Render(appContext.Grid)
			case "<MouseLeft>":
				payload := e.Payload.(ui.Mouse)
				appContext.View.Workplaces.HandleClick(payload.X, payload.Y)
				ui.Render(appContext.View.Workplaces)
			case "k", "<Up>", "<MouseWheelUp>":
				appContext.View.Workplaces.ScrollUp()
				ui.Render(appContext.View.Workplaces)
			case "j", "<Down>", "<MouseWheelDown>":
				appContext.View.Workplaces.ScrollDown()
				ui.Render(appContext.View.Workplaces)
			case "<Home>":
				appContext.View.Workplaces.ScrollTop()
				ui.Render(appContext.View.Workplaces)
			case "g":
				if previousKey == "g" {
					appContext.View.Workplaces.ScrollTop()
					ui.Render(appContext.View.Workplaces)
				}
			case "G", "<End>":
				appContext.View.Workplaces.ScrollBottom()
				ui.Render(appContext.View.Workplaces)
			case "<C-d>":
				appContext.View.Workplaces.ScrollHalfPageDown()
				ui.Render(appContext.View.Workplaces)
			case "<C-u>":
				appContext.View.Workplaces.ScrollHalfPageUp()
				ui.Render(appContext.View.Workplaces)
			case "<C-f>":
				appContext.View.Workplaces.ScrollPageDown()
				ui.Render(appContext.View.Workplaces)
			case "<C-b>":
				appContext.View.Workplaces.ScrollPageUp()
				ui.Render(appContext.View.Workplaces)
			case "<Enter>":
				appContext.View.Workplaces.SelectWorkplace()
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
		log.Fatalf("failed to create AppContext: %v", err)
	}

	eventLoop(appContext)
}
