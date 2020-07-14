package views

import (
	"github.com/abgeo/goclockify/src/config"
	"github.com/abgeo/goclockify/src/services"
	cw "github.com/abgeo/goclockify/src/widgets"
	w "github.com/abgeo/goclockify/src/widgets"
)

type View struct {
	Config      *config.Config
	User        *w.UserWidget
	Workplaces  *w.WorkplacesWidget
	TimeEntries *w.TimeEntriesWidget
}

func CreateView(config *config.Config, clockifyService *services.ClockifyService) (*View, error) {
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

	// Setup TimeEntriesWidget

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

	return &View{
		Config:      config,
		User:        user,
		Workplaces:  workplaces,
		TimeEntries: timeEntries,
	}, nil
}
