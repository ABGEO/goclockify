package widgets

import (
	"fmt"
	cui "github.com/abgeo/goclockify/src/goclockify"
	"github.com/abgeo/goclockify/src/services"
	ui "github.com/gizak/termui/v3"
)

type WorkplacesWidget struct {
	*cui.Table
	Workplaces []services.Workplace
}

func NewWorkplacesWidget(clockifyService *services.ClockifyService) (*WorkplacesWidget, error) {
	self := &WorkplacesWidget{
		Table: cui.NewTable(),
	}

	workplaces, err := clockifyService.GetWorkplaces()
	if err != nil {
		return nil, err
	}

	self.Title = " Workplaces "
	self.ShowCursor = true
	self.ShowLocation = true
	self.ColGap = 3
	self.PadLeft = 2
	self.UniqueCol = 0
	self.Header = []string{"ID", "Name"}
	self.Workplaces = workplaces
	self.ColResizer = func() {
		self.ColWidths = []int{5, maxInt(self.Inner.Dx()-26, 10)}
	}

	self.workplacesToRows()

	return self, nil
}

func (self *WorkplacesWidget) workplacesToRows() {
	var workplaces *[]services.Workplace
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
