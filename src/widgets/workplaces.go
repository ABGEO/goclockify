package widgets

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"strconv"
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
	self.PadLeft = 3
	self.Header = []string{"", "Name"}
	self.ColResizer = func() {
		self.ColWidths = []int{0, ui.MinInt(self.Inner.Dx()-26, 50)}
	}

	return self
}

func (self *WorkplacesWidget) SetWorkplaces(workplaces []Workplace) {
	self.Workplaces = workplaces
	self.workplacesToRows()
}

func (self *WorkplacesWidget) GetSelectedWorkplace() (Workplace, error) {
	selectedIndex := self.Rows[self.SelectedRow][0]
	i, err := strconv.Atoi(selectedIndex)
	if err != nil {
		return Workplace{}, err
	}

	return self.Workplaces[i], nil
}

func (self *WorkplacesWidget) workplacesToRows() {
	var workplaces *[]Workplace
	workplaces = &self.Workplaces
	strings := make([][]string, len(*workplaces))
	for i, w := range *workplaces {
		strings[i] = make([]string, 2)
		strings[i][0] = fmt.Sprintf("%d", i)
		strings[i][1] = w.Name
	}
	self.Rows = strings
}
