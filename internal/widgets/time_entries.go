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
	"time"
)

// Tag represents the tag entity from the API
type Tag struct {
	ID   string
	Name string
}

// Project represents the project entity from the API
type Project struct {
	ID         string
	Name       string
	ClientName string
}

// TimeInterval represents the time interval entity from the API
type TimeInterval struct {
	Start    time.Time
	End      time.Time
	Duration string
}

// TimeEntry represents the time entry entity from the API
type TimeEntry struct {
	ID           string
	Description  string
	Tags         []Tag
	Billable     bool
	Project      Project
	TimeInterval TimeInterval
}

// TimeEntriesWidget is a component with the time entries
type TimeEntriesWidget struct {
	*components.Table
	TimeEntries []TimeEntry
}

// NewTimeEntriesWidget creates new TimeEntriesWidget
func NewTimeEntriesWidget() *TimeEntriesWidget {
	self := &TimeEntriesWidget{
		Table: components.NewTable(),
	}

	self.ShowCursor = true
	self.ShowLocation = true
	self.ColGap = 3
	self.Header = []string{"", "Description", "Duration", ""}
	self.ColResizer = func() {
		self.ColWidths = []int{0, ui.MinInt(self.Inner.Dx()-26, 100), 15, 1}
	}

	return self
}

// UpdateData updates and reloads TimeEntriesWidget
func (w *TimeEntriesWidget) UpdateData(timeEntries []TimeEntry, workplace Workplace) {
	w.TimeEntries = timeEntries
	w.Title = fmt.Sprintf(" %s - Time Entries ", workplace.Name)
	w.SelectedItem = ""
	w.entriesToRows()
	w.ScrollTop()
}

// GetSelectedTimeEntry returns the selected time entry
func (w *TimeEntriesWidget) GetSelectedTimeEntry() (TimeEntry, error) {
	selectedIndex := w.Rows[w.SelectedRow][0]
	i, err := strconv.Atoi(selectedIndex)
	if err != nil {
		return TimeEntry{}, err
	}

	return w.TimeEntries[i], nil
}

func (w *TimeEntriesWidget) entriesToRows() {
	var timeEntries *[]TimeEntry
	timeEntries = &w.TimeEntries
	strings := make([][]string, len(*timeEntries))
	for i, t := range *timeEntries {
		strings[i] = make([]string, 4)
		strings[i][0] = fmt.Sprintf("%d", i)
		strings[i][1] = t.Description

		if t.TimeInterval.End.IsZero() {
			strings[i][2] = "Running..."
		} else {
			strings[i][2] = t.TimeInterval.End.Sub(t.TimeInterval.Start).String()
		}

		if t.Billable {
			strings[i][3] = "$"
		} else {
			strings[i][3] = ""
		}
	}
	w.Rows = strings
}
