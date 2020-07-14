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
	self.ColGap = 3
	self.PadLeft = 2
	self.UniqueCol = 0
	self.Header = []string{"No", "Name"}
	self.ColResizer = func() {
		self.ColWidths = []int{2, ui.MaxInt(self.Inner.Dx()-26, 10)}
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

	return self.Workplaces[i-1], nil
}

func (self *WorkplacesWidget) workplacesToRows() {
	var workplaces *[]Workplace
	workplaces = &self.Workplaces
	strings := make([][]string, len(*workplaces))
	for i, w := range *workplaces {
		strings[i] = make([]string, 2)
		strings[i][0] = fmt.Sprintf("%d", i+1)
		strings[i][1] = w.Name
	}
	self.Rows = strings
}

func (self *WorkplacesWidget) SelectWorkplace() {
	self.SelectedItem = ""

	ui.Clear()
	fmt.Println(self.GetSelectedWorkplace())
}
