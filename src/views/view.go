package views

import (
	cw "github.com/abgeo/goclockify/src/widgets"
	w "github.com/abgeo/goclockify/src/widgets"
)

type View struct {
	Workplaces *w.WorkplacesWidget
}

func CreateView() (*View, error) {
	workplaces, err := cw.NewWorkplacesWidget()
	if err != nil {
		return nil, err
	}

	return &View{
		Workplaces: workplaces,
	}, nil
}
