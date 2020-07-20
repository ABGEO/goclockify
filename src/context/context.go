// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package context

import (
	"github.com/abgeo/goclockify/src/config"
	"github.com/abgeo/goclockify/src/services"
	"github.com/abgeo/goclockify/src/theme"
	"github.com/abgeo/goclockify/src/views"
	ui "github.com/gizak/termui/v3"
)

// AppContext stores application context
type AppContext struct {
	Grid            *ui.Grid
	View            *views.View
	Config          *config.Config
	ClockifyService *services.ClockifyService
	Theme           theme.Theme
}

// CreateAppContext creates new application context object
func CreateAppContext() (*AppContext, error) {
	context := &AppContext{
		Theme: theme.Default,
	}

	var err error
	context.Config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	context.ClockifyService, err = services.NewClockifyService(context.Config)
	if err != nil {
		return nil, err
	}

	ui.Theme.Default = ui.NewStyle(ui.Color(context.Theme.Fg), ui.Color(context.Theme.Bg))
	ui.Theme.Block.Title = ui.NewStyle(ui.Color(context.Theme.BorderLabel), ui.Color(context.Theme.Bg))
	ui.Theme.Block.Border = ui.NewStyle(ui.Color(context.Theme.BorderLine), ui.Color(context.Theme.Bg))

	context.View, err = views.CreateView(context.Config, context.ClockifyService)
	if err != nil {
		return nil, err
	}

	context.View.Workplaces.CursorColor = ui.Color(context.Theme.TableCursor)
	context.View.TimeEntries.CursorColor = ui.Color(context.Theme.TableCursor)

	context.Grid = ui.NewGrid()
	context.Grid.Set(
		ui.NewCol(1.0/4,
			ui.NewRow(1.0/3, context.View.User),
			ui.NewRow(2.0/3, context.View.Workplaces),
		),
		ui.NewCol(3.0/4, context.View.TimeEntries),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	context.Grid.SetRect(0, 0, termWidth, termHeight)

	ui.Render(context.Grid)

	context.View.TimeEntries.ScrollTop()
	ui.Render(context.View.TimeEntries)

	return context, nil
}
