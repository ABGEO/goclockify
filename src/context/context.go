package context

import (
	ColorSchemes "github.com/abgeo/goclockify/colorschemes"
	"github.com/abgeo/goclockify/src/views"
	ui "github.com/gizak/termui/v3"
)

var (
	ColorScheme = ColorSchemes.Default
)

type AppContext struct {
	Grid *ui.Grid
	View *views.View
}

func CreateAppContext() (*AppContext, error) {
	ui.Theme.Default = ui.NewStyle(ui.Color(ColorScheme.Fg), ui.Color(ColorScheme.Bg))
	ui.Theme.Block.Title = ui.NewStyle(ui.Color(ColorScheme.BorderLabel), ui.Color(ColorScheme.Bg))
	ui.Theme.Block.Border = ui.NewStyle(ui.Color(ColorScheme.BorderLine), ui.Color(ColorScheme.Bg))

	view, err := views.CreateView()
	if err != nil {
		return nil, err
	}

	view.Workplaces.CursorColor = ui.Color(ColorScheme.TableCursor)

	grid := ui.NewGrid()
	grid.Set(
		ui.NewCol(1.0/4, view.Workplaces),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	ui.Render(grid)

	return &AppContext{
		Grid: grid,
		View: view,
	}, nil
}
