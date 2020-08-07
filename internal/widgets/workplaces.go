// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package widgets

import (
	"fmt"
	"github.com/abgeo/goclockify/internal/components"
	ui "github.com/gizak/termui/v3"
	"strconv"
)

// Workplace represents the workplace entity from the API
type Workplace struct {
	ID   string
	Name string
}

// WorkplacesWidget is a component that displays workplaces
type WorkplacesWidget struct {
	*components.Table
	Workplaces []Workplace
}

// NewWorkplacesWidget creates new WorkplacesWidget
func NewWorkplacesWidget() *WorkplacesWidget {
	self := &WorkplacesWidget{
		Table: components.NewTable(),
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

// SetWorkplaces sets the value of WorkplacesWidget.Workplaces
func (w *WorkplacesWidget) SetWorkplaces(workplaces []Workplace) {
	w.Workplaces = workplaces
	w.workplacesToRows()
}

// GetSelectedWorkplace returns the selected workplace
func (w *WorkplacesWidget) GetSelectedWorkplace() (Workplace, error) {
	selectedIndex := w.Rows[w.SelectedRow][0]
	i, err := strconv.Atoi(selectedIndex)
	if err != nil {
		return Workplace{}, err
	}

	return w.Workplaces[i], nil
}

func (w *WorkplacesWidget) workplacesToRows() {
	var workplaces *[]Workplace
	workplaces = &w.Workplaces
	strings := make([][]string, len(*workplaces))
	for i, w := range *workplaces {
		strings[i] = make([]string, 2)
		strings[i][0] = fmt.Sprintf("%d", i)
		strings[i][1] = w.Name
	}
	w.Rows = strings
}
