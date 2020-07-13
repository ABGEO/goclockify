package views

import (
	"github.com/abgeo/goclockify/src/config"
	"github.com/abgeo/goclockify/src/services"
	cw "github.com/abgeo/goclockify/src/widgets"
	w "github.com/abgeo/goclockify/src/widgets"
)

type View struct {
	Config     *config.Config
	Workplaces *w.WorkplacesWidget
}

func CreateView(config *config.Config, clockifyService *services.ClockifyService) (*View, error) {
	workplaces, err := cw.NewWorkplacesWidget()
	if err != nil {
		return nil, err
	}

	workplaceItems, err := clockifyService.GetWorkplaces()
	if err != nil {
		return nil, err
	}

	workplaces.SetWorkplaces(workplaceItems)

	return &View{
		Config:     config,
		Workplaces: workplaces,
	}, nil
}
