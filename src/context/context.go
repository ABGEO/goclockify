package context

import (
	ColorSchemes "github.com/abgeo/goclockify/colorschemes"
	"github.com/abgeo/goclockify/src/config"
	"github.com/abgeo/goclockify/src/services"
	"github.com/abgeo/goclockify/src/views"
	ui "github.com/gizak/termui/v3"
)

type AppContext struct {
	Grid            *ui.Grid
	View            *views.View
	Config          *config.Config
	ClockifyService *services.ClockifyService
	ColorScheme     ColorSchemes.ColorScheme
}

func CreateAppContext() (*AppContext, error) {
	context := &AppContext{
		ColorScheme: ColorSchemes.Default,
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

	ui.Theme.Default = ui.NewStyle(ui.Color(context.ColorScheme.Fg), ui.Color(context.ColorScheme.Bg))
	ui.Theme.Block.Title = ui.NewStyle(ui.Color(context.ColorScheme.BorderLabel), ui.Color(context.ColorScheme.Bg))
	ui.Theme.Block.Border = ui.NewStyle(ui.Color(context.ColorScheme.BorderLine), ui.Color(context.ColorScheme.Bg))

	context.View, err = views.CreateView(context.Config, context.ClockifyService)
	if err != nil {
		return nil, err
	}

	context.View.Workplaces.CursorColor = ui.Color(context.ColorScheme.TableCursor)
	context.View.TimeEntries.CursorColor = ui.Color(context.ColorScheme.TableCursor)

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
