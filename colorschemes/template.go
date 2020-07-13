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
