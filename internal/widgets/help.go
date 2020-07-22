// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package widgets

import (
	w "github.com/gizak/termui/v3/widgets"
)

const help = `
Quit: q and <C-c>

Workspaces navigation
  - a: up
  - z: down

Time Entries navigation
  - k, <Up> and <MouseWheelUp>: up
  - j, <Down> and <MouseWheelDown>: down
  - g and <Home>: jump to top
  - G and <End>: jump to bottom
  - <Enter> display time entry details

Other
  - <Escape>: close the 2nd level window, go to the dashboard
  - <F1> and ? show this message
`

// HelpWidget is a component with the help text
type HelpWidget struct {
	w.Paragraph
}

// NewHelpWidget creates new HelpWidget
func NewHelpWidget() *HelpWidget {
	self := &HelpWidget{
		Paragraph: *w.NewParagraph(),
	}

	self.Title = " Help "
	self.Text = help

	return self
}
