// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package ColorSchemes

/*
	The standard 256 terminal colors are supported.

	-1 = clear

	You can combine a color with 'Bold', 'Underline', or 'Reverse' by using bitwise OR ('|') and the name of the Color.
	For example, to get Bold red Labels, you would do 'Labels: 2 | Bold'.
*/

const (
	Bold int = 1 << (iota + 9)
	Underline
	Reverse
)

type ColorScheme struct {
	Name   string
	Author string

	Fg int
	Bg int

	BorderLabel int
	BorderLine  int

	TableCursor int
}
