package widgets

import (
	w "github.com/gizak/termui/v3/widgets"
)

const help = `
Quit: q or <C-c>

Workspaces navigation
  - a: up
  - z: down

Time Entries navigation
  - k, <Up> and <MouseWheelUp>: up
  - j, <Down> and <MouseWheelDown>: down
  - g and <Home>: jump to top
  - G and <End>: jump to bottom
  - <Enter> display time entry details
  - <F1> and ? show this message
`

type HelpWidget struct {
	w.Paragraph
}

func NewHelpWidget() *HelpWidget {
	self := &HelpWidget{
		Paragraph: *w.NewParagraph(),
	}

	self.Title = " Help "
	self.Text = help

	return self
}
