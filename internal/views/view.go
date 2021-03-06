// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package views

import (
	"github.com/abgeo/goclockify/configs"
	"github.com/abgeo/goclockify/internal/services"
	cw "github.com/abgeo/goclockify/internal/widgets"
	w "github.com/abgeo/goclockify/internal/widgets"
)

// View stores global widget instances
type View struct {
	Config      *configs.Config
	User        *w.UserWidget
	Workplaces  *w.WorkplacesWidget
	TimeEntries *w.TimeEntriesWidget
	TimeEntry   *w.TimeEntryWidget
	Help        *w.HelpWidget
}

// CreateView creates a new View instance
func CreateView(config *configs.Config, clockifyService *services.ClockifyService) (*View, error) {
	// Setup UserWidget.

	user := cw.NewUserWidget()
	user.SetUser(clockifyService.CurrentUser)

	// Setup WorkplacesWidget.

	workplaces := cw.NewWorkplacesWidget()

	workplaceItems, err := clockifyService.GetWorkplaces()
	if err != nil {
		return nil, err
	}

	workplaces.SetWorkplaces(workplaceItems)

	// Setup TimeEntriesWidget.

	timeEntries := w.NewTimeEntriesWidget()

	selectedWorkspace, err := workplaces.GetSelectedWorkplace()
	if err != nil {
		return nil, err
	}

	timeEntryItems, err := clockifyService.GetTimeEntries(selectedWorkspace.ID)
	if err != nil {
		return nil, err
	}

	timeEntries.UpdateData(timeEntryItems, selectedWorkspace)

	// Setup TimeEntryWidget.
	timeEntry := w.NewTimeEntryWidget()

	// Setup HelpWidget.
	help := w.NewHelpWidget()

	return &View{
		Config:      config,
		User:        user,
		Workplaces:  workplaces,
		TimeEntries: timeEntries,
		TimeEntry:   timeEntry,
		Help:        help,
	}, nil
}
