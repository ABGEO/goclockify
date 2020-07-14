package context

import (
	ColorSchemes "github.com/abgeo/goclockify/colorschemes"
	"github.com/abgeo/goclockify/src/config"
	"github.com/abgeo/goclockify/src/services"
	"github.com/abgeo/goclockify/src/views"
	ui "github.com/gizak/termui/v3"
)

var (
	ColorScheme = ColorSchemes.Default
)

type AppContext struct {
	Grid   *ui.Grid
	View   *views.View
	Config *config.Config
}

func CreateAppContext() (*AppContext, error) {
	appConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	clockifyService, err := services.NewClockifyService(appConfig)
	if err != nil {
		return nil, err
	}

	ui.Theme.Default = ui.NewStyle(ui.Color(ColorScheme.Fg), ui.Color(ColorScheme.Bg))
	ui.Theme.Block.Title = ui.NewStyle(ui.Color(ColorScheme.BorderLabel), ui.Color(ColorScheme.Bg))
	ui.Theme.Block.Border = ui.NewStyle(ui.Color(ColorScheme.BorderLine), ui.Color(ColorScheme.Bg))

	view, err := views.CreateView(appConfig, clockifyService)
	if err != nil {
		return nil, err
	}

	view.Workplaces.CursorColor = ui.Color(ColorScheme.TableCursor)
	view.TimeEntries.CursorColor = ui.Color(ColorScheme.TableCursor)

	grid := ui.NewGrid()
	grid.Set(
		ui.NewCol(1.0/4,
			ui.NewRow(1.0/3, view.User),
			ui.NewRow(2.0/3, view.Workplaces),
		),
		ui.NewCol(3.0/4, view.TimeEntries),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	ui.Render(grid)

	return &AppContext{
		Grid:   grid,
		View:   view,
		Config: appConfig,
	}, nil
}
