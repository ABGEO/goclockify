// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package components

import (
	ui "github.com/gizak/termui/v3"
	"github.com/nsf/termbox-go"
	"image"
)

// Input widget type
type Input struct {
	ui.Block
	text          string
	TextStyle     ui.Style
	ActionMapping map[string]func()
	cursorX       int
	isCapturing   bool
}

// NewInput creates new Input widget
func NewInput() *Input {
	return &Input{
		Block:     *ui.NewBlock(),
		text:      "",
		TextStyle: ui.Theme.Paragraph.Text,

		cursorX: 0,
	}
}

// Capture begins catching events and updates the content of the text field.
func (i *Input) Capture() {
	i.isCapturing = true

	ui.Render(i)
EventLoop:
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<Left>":
				i.moveLeft()
			case "<Right>":
				i.moveRight()
			case "<Backspace>":
				i.removeLeft()
			case "<Delete>":
				i.removeRight()
			case "<Escape>":
				break EventLoop
			default:
				if callback, ok := i.ActionMapping[e.ID]; ok {
					callback()
					break EventLoop
				} else if char := i.getChar(e.ID); char != byte(0) {
					i.AddText(string(char))
				} else {
					break EventLoop
				}
			}

			ui.Render(i)
		}
	}

	i.isCapturing = false
	ui.Render(i)
}

// AddText updates the content of the text field based on cursor position.
func (i *Input) AddText(s string) {
	if i.cursorX == 0 {
		i.text = s + i.text
	} else if i.cursorX == len(i.text) {
		i.text += s
	} else {
		i.text = i.text[:i.cursorX] + s + i.text[i.cursorX:]
	}

	i.cursorX += len(s)
}

// GetText returns the content of the text field.
func (i *Input) GetText() string {
	return i.text
}

func (i *Input) removeLeft() {
	if i.cursorX > 0 {
		i.text = i.text[:i.cursorX-1] + i.text[i.cursorX:]

		i.cursorX--
	}
}

func (i *Input) moveLeft() {
	if i.cursorX > 0 {
		i.cursorX--
	}
}

func (i *Input) moveRight() {
	if i.cursorX < len(i.text) {
		i.cursorX++
	}
}

func (i *Input) removeRight() {
	if i.cursorX < len(i.text) {
		i.text = i.text[:i.cursorX] + i.text[i.cursorX+1:]
	}
}

func (i *Input) getChar(s string) byte {
	mapping := map[string]byte{
		"<Space>": ' ',
	}

	if val, ok := mapping[s]; ok {
		return val
	}

	if len(s) == 1 {
		return s[0]
	}

	return byte(0)
}

// Draw implements the ui.Drawable interface.
func (i *Input) Draw(buf *ui.Buffer) {
	i.Block.Draw(buf)

	cells := ui.ParseStyles(i.text, i.TextStyle)
	rows := ui.SplitCells(cells, '\n')

	for y, row := range rows {
		if y+i.Inner.Min.Y >= i.Inner.Max.Y {
			break
		}

		row = ui.TrimCells(row, i.Inner.Dx())
		for _, cx := range ui.BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(i.Inner.Min))
		}
	}

	if i.isCapturing {
		cursorXOffset := i.Min.X
		cursorYOffset := i.Min.Y

		if i.BorderTop {
			cursorYOffset++
		}

		if i.BorderLeft {
			cursorXOffset++
		}

		termbox.SetCursor(cursorXOffset+i.cursorX, cursorYOffset)
	} else {
		termbox.HideCursor()
	}
}
