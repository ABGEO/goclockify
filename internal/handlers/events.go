// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package handlers

import (
	"github.com/abgeo/goclockify/internal/components"
	"github.com/abgeo/goclockify/internal/context"
	ui "github.com/gizak/termui/v3"
	"image"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	updateInterval = time.Second

	showDashboard       = true
	showSingleTimeEntry = false
	showHelp            = false

	blockMainInput = false

	shownConfirm *components.Confirm
)

// Initialize new events handler
func Initialize(appContext *context.AppContext) {
	drawTicker := time.NewTicker(updateInterval).C

	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)

	uiEvents := ui.PollEvents()
	actionMap := createActionMap(appContext)

	for {
		select {
		case <-sigTerm:
			actionQuit(nil, nil)
		case <-drawTicker:
			actionDrawTicker(appContext)
		case e := <-uiEvents:
			switch e.Type {
			case ui.ResizeEvent:
				actionResize(appContext, &e)
			case ui.KeyboardEvent, ui.MouseEvent:
				if !blockMainInput {
					if action, ok := actionMap[e.ID]; ok {
						action(appContext, &e)
					}
				}

				switch e.ID {
				case "<Escape>":
					if shownConfirm != nil {
						blockMainInput = false
						shownConfirm = nil
						conditionalRender(showDashboard, appContext.Grid)
					}
				case "<Enter>":
					if shownConfirm != nil {
						shownConfirm.Callback()
						blockMainInput = false
						shownConfirm = nil
						conditionalRender(showDashboard, appContext.Grid)
					}
				}
			}
		}
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

func createActionMap(appContext *context.AppContext) (mapping map[string]func(*context.AppContext, *ui.Event)) {
	mapping = make(map[string]func(*context.AppContext, *ui.Event))
	ckm := appContext.Config.KeyMapping

	// Workspace mapping

	for _, k := range ckm.Workspace.NavigationUp {
		mapping[k] = actionWorkspaceNavUp
	}

	for _, k := range ckm.Workspace.NavigationDown {
		mapping[k] = actionWorkspaceNavDown
	}

	// Time Entries mapping

	for _, k := range ckm.TimeEntries.NavigationUp {
		mapping[k] = actionTimeEntriesNavUp
	}

	for _, k := range ckm.TimeEntries.NavigationDown {
		mapping[k] = actionTimeEntriesNavDown
	}

	for _, k := range ckm.TimeEntries.NavigationToBottom {
		mapping[k] = actionTimeEntriesNavToBottom
	}

	for _, k := range ckm.TimeEntries.NavigationToTop {
		mapping[k] = actionTimeEntriesNavToTop
	}

	for _, k := range ckm.TimeEntries.NavigationSelect {
		mapping[k] = actionTimeEntriesNavSelect
	}

	for _, k := range ckm.TimeEntries.Add {
		mapping[k] = actionTimeEntriesAdd
	}

	for _, k := range ckm.TimeEntries.Delete {
		mapping[k] = actionTimeEntriesDelete
	}

	for _, k := range ckm.TimeEntries.Edit {
		mapping[k] = actionTimeEntriesEdit
	}

	// Other mapping

	for _, k := range ckm.Other.Quit {
		mapping[k] = actionQuit
	}

	for _, k := range ckm.Other.Help {
		mapping[k] = actionHelp
	}

	for _, k := range ckm.Other.CloseWindow {
		mapping[k] = actionCloseWindow
	}

	// Mouse mapping
	mapping["<MouseLeft>"] = actionMouseLeft

	return mapping
}

func actionDrawTicker(appContext *context.AppContext) {
	if showSingleTimeEntry {
		appContext.View.TimeEntry.UpdateTable()
		ui.Render(appContext.View.TimeEntry)
	}

	conditionalRender(showDashboard, appContext.Grid)

	if shownConfirm != nil {
		shownConfirm.Render()
	}
}

func actionResize(appContext *context.AppContext, e *ui.Event) {
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

	if shownConfirm != nil {
		shownConfirm.Location = image.Point{
			X: payload.Width / 2,
			Y: payload.Height / 2,
		}
		shownConfirm.Render()
	}
}

func actionTimeEntriesNavUp(appContext *context.AppContext, _ *ui.Event) {
	appContext.View.TimeEntries.ScrollUp()
	conditionalRender(showDashboard, appContext.View.TimeEntries)
}

func actionTimeEntriesNavDown(appContext *context.AppContext, _ *ui.Event) {
	appContext.View.TimeEntries.ScrollDown()
	conditionalRender(showDashboard, appContext.View.TimeEntries)
}

func actionTimeEntriesNavSelect(appContext *context.AppContext, _ *ui.Event) {
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

func actionTimeEntriesAdd(appContext *context.AppContext, e *ui.Event) {
	blockMainInput = true
	showDashboard = false

	form := components.NewTimeEntryForm()
	form.SubmitCallback = func() {
		workplace, _ := appContext.View.Workplaces.GetSelectedWorkplace()
		err := appContext.ClockifyService.AddTimeEntry(workplace.ID, form.GetData())
		if err == nil {
			updateTimeEntries(appContext)
		}

		blockMainInput = false
		actionCloseWindow(appContext, e)
	}

	ui.Clear()
	form.Render()

Loop:
	for {
		event := <-ui.PollEvents()
		switch event.ID {
		case "<Tab>":
			form.ActiveInput--
			form.FocusNext()
		case "<Enter>":
			form.SubmitCallback()
			break Loop
		case "<Escape>":
			actionCloseWindow(appContext, e)
			blockMainInput = false
			break Loop
		}
	}
}

func actionTimeEntriesDelete(appContext *context.AppContext, _ *ui.Event) {
	blockMainInput = true
	shownConfirm = components.NewConfirm()
	shownConfirm.Text = "Do you want to remove this Time Entry?"

	terminalWidth, terminalHeight := ui.TerminalDimensions()
	shownConfirm.Location = image.Point{
		X: terminalWidth / 2,
		Y: terminalHeight / 2,
	}

	shownConfirm.Callback = func() {
		workplace, _ := appContext.View.Workplaces.GetSelectedWorkplace()
		timeEntry, _ := appContext.View.TimeEntries.GetSelectedTimeEntry()

		err := appContext.ClockifyService.DeleteTimeEntry(workplace.ID, timeEntry.ID)
		if err == nil {
			updateTimeEntries(appContext)
			conditionalRender(showDashboard, appContext.View.TimeEntries)
		}
	}

	shownConfirm.Render()
}

func actionTimeEntriesEdit(appContext *context.AppContext, e *ui.Event) {
	blockMainInput = true
	showDashboard = false

	form := components.NewTimeEntryForm()
	form.Data, _ = appContext.View.TimeEntries.GetSelectedTimeEntry()
	form.SubmitCallback = func() {
		workplace, _ := appContext.View.Workplaces.GetSelectedWorkplace()
		err := appContext.ClockifyService.EditTimeEntry(workplace.ID, form.GetData())
		if err == nil {
			updateTimeEntries(appContext)
		}

		blockMainInput = false
		actionCloseWindow(appContext, e)
	}

	ui.Clear()
	form.Render()

Loop:
	for {
		event := <-ui.PollEvents()
		switch event.ID {
		case "<Tab>":
			form.ActiveInput--
			form.FocusNext()
		case "<Enter>":
			form.SubmitCallback()
			break Loop
		case "<Escape>":
			actionCloseWindow(appContext, e)
			blockMainInput = false
			break Loop
		}
	}
}

func actionTimeEntriesNavToBottom(appContext *context.AppContext, _ *ui.Event) {
	appContext.View.TimeEntries.ScrollBottom()
	conditionalRender(showDashboard, appContext.View.TimeEntries)
}

func actionTimeEntriesNavToTop(appContext *context.AppContext, _ *ui.Event) {
	appContext.View.TimeEntries.ScrollTop()
	conditionalRender(showDashboard, appContext.View.TimeEntries)
}

func actionWorkspaceNavDown(appContext *context.AppContext, _ *ui.Event) {
	appContext.View.Workplaces.ScrollDown()
	updateTimeEntries(appContext)
	conditionalRender(showDashboard, appContext.View.Workplaces)
	conditionalRender(showDashboard, appContext.View.TimeEntries)
}

func actionWorkspaceNavUp(appContext *context.AppContext, _ *ui.Event) {
	appContext.View.Workplaces.ScrollUp()
	updateTimeEntries(appContext)
	conditionalRender(showDashboard, appContext.View.Workplaces)
	conditionalRender(showDashboard, appContext.View.TimeEntries)
}

func actionQuit(_ *context.AppContext, _ *ui.Event) {
	ui.Close()
	os.Exit(0)
}

func actionHelp(appContext *context.AppContext, _ *ui.Event) {
	showHelp = !showHelp
	showDashboard = !showHelp

	if showHelp {
		ui.Clear()
		terminalWidth, terminalHeight := ui.TerminalDimensions()
		appContext.View.Help.SetRect(0, 0, terminalWidth, terminalHeight)
		ui.Render(appContext.View.Help)
	}
}

func actionCloseWindow(appContext *context.AppContext, _ *ui.Event) {
	if !showDashboard {
		showDashboard = true
		showSingleTimeEntry = false
		showHelp = false
		ui.Clear()
		ui.Render(appContext.Grid)
	}
}

func actionMouseLeft(appContext *context.AppContext, e *ui.Event) {
	payload := e.Payload.(ui.Mouse)
	appContext.View.TimeEntries.HandleClick(payload.X, payload.Y)
	conditionalRender(showDashboard, appContext.View.TimeEntries)
}
