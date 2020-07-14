package widgets

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
)

type Workplace struct {
	ID   string
	Name string
}

type WorkplacesWidget struct {
	*Table
	Workplaces []Workplace
}

func NewWorkplacesWidget() *WorkplacesWidget {
	self := &WorkplacesWidget{
		Table: NewTable(),
	}

	self.Title = " Workplaces "
	self.ShowCursor = true
	self.ShowLocation = true
	self.ColGap = 3
	self.PadLeft = 2
	self.UniqueCol = 0
	self.Header = []string{"ID", "Name"}
	self.ColResizer = func() {
		self.ColWidths = []int{5, maxInt(self.Inner.Dx()-26, 10)}
	}

	return self
}

func (self *WorkplacesWidget) SetWorkplaces(workplaces []Workplace) {
	self.Workplaces = workplaces
	self.workplacesToRows()
}

func (self *WorkplacesWidget) workplacesToRows() {
	var workplaces *[]Workplace
	workplaces = &self.Workplaces
	strings := make([][]string, len(*workplaces))
	for i := range *workplaces {
		strings[i] = make([]string, 2)
		strings[i][0] = (*workplaces)[i].ID
		strings[i][1] = (*workplaces)[i].Name
	}
	self.Rows = strings
}

func (self *WorkplacesWidget) SelectWorkplace() {
	self.SelectedItem = ""
	selectedWorkplace := self.Rows[self.SelectedRow][self.UniqueCol]

	ui.Clear()
	fmt.Println(selectedWorkplace)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
