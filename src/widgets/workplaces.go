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

	self.Title = " Workspaces "
	self.ShowCursor = true
	self.ShowLocation = true
	self.ColGap = 3
	self.PadLeft = 2
	self.UniqueCol = 0
	self.Header = []string{"ID", "Name"}
	self.ColResizer = func() {
		self.ColWidths = []int{5, ui.MaxInt(self.Inner.Dx()-26, 10)}
	}

	return self
}

func (self *WorkplacesWidget) SetWorkplaces(workplaces []Workplace) {
	self.Workplaces = workplaces
	self.workplacesToRows()
}

func (self *WorkplacesWidget) GetSelectedWorkplace() string {
	return self.Rows[self.SelectedRow][self.UniqueCol]
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

	ui.Clear()
	fmt.Println(self.GetSelectedWorkplace())
}
