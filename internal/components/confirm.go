// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package components

import (
	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
	"image"
)

// Confirm widget type
type Confirm struct {
	P0, P1, P2, P3 *w.Paragraph

	Text string

	Location image.Point

	Callback func()
}

// NewConfirm creates new Confirm widget
func NewConfirm() *Confirm {
	c := &Confirm{
		P0:       w.NewParagraph(),
		P1:       w.NewParagraph(),
		P2:       w.NewParagraph(),
		P3:       w.NewParagraph(),
		Callback: func() {},
	}

	c.P2.Text = "No <Escape>"
	c.P2.TextStyle = ui.Style{
		Fg: ui.ColorClear,
		Bg: ui.ColorYellow,
	}

	c.P3.Text = "Yes <Enter>"
	c.P3.TextStyle = ui.Style{
		Fg: ui.ColorClear,
		Bg: ui.ColorRed,
	}

	return c
}

// Render renders Confirm widget
func (c *Confirm) Render() {
	x := c.Location.X
	y := c.Location.Y

	c.P1.Text = c.Text

	c.P0.SetRect(x-26, y-4, x+26, y+5)
	c.P1.SetRect(x-25, y-3, x+25, y+1)
	c.P2.SetRect(x-25, y+1, x-12, y+4)
	c.P3.SetRect(x-12, y+1, x+1, y+4)

	ui.Render(c.P0, c.P1, c.P2, c.P3)
}
