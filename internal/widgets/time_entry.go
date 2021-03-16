// This file is part of the abgeo/goclokify.
//
// Copyright (C) 2020 Temuri Takalandze <takalandzet@gmail.com>.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package widgets

import (
	"github.com/abgeo/goclockify/internal/types"
	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
	"time"
)

// TimeEntryWidget is a component with the single time entry data
type TimeEntryWidget struct {
	*w.Table
	TimeEntry types.TimeEntry
}

// NewTimeEntryWidget creates new TimeEntryWidget
func NewTimeEntryWidget() *TimeEntryWidget {
	self := &TimeEntryWidget{
		Table: w.NewTable(),
	}

	self.Title = " Time Entry Details "
	self.FillRow = true
	terminalWidth, _ := ui.TerminalDimensions()
	self.ColumnWidths = []int{15, terminalWidth - 15}
	self.Rows = [][]string{
		{"Description", ""},
		{"Project", ""},
		{"Client", ""},
		{"Tags", ""},
		{"Start", ""},
		{"End", ""},
		{"Duration", ""},
	}

	return self
}

// SetTimeEntry sets the value of TimeEntryWidget.TimeEntry
func (t *TimeEntryWidget) SetTimeEntry(timeEntry types.TimeEntry) {
	t.TimeEntry = timeEntry
	t.UpdateTable()
}

// UpdateTable updates table with TimeEntryWidget.TimeEntry data
func (t *TimeEntryWidget) UpdateTable() {
	t.Rows[0][1] = t.TimeEntry.Description
	t.Rows[1][1] = t.TimeEntry.Project.Name
	t.Rows[2][1] = t.TimeEntry.Project.ClientName
	t.Rows[4][1] = t.TimeEntry.TimeInterval.Start.Format("01/02/2006 15:04:05")
	t.Rows[5][1] = t.TimeEntry.TimeInterval.End.Format("01/02/2006 15:04:05")

	if t.TimeEntry.TimeInterval.End.IsZero() {
		t.Rows[6][1] = "Running - " + time.Now().Sub(t.TimeEntry.TimeInterval.Start).String()
	} else {
		t.Rows[6][1] = t.TimeEntry.TimeInterval.End.Sub(t.TimeEntry.TimeInterval.Start).String()
	}

	t.Rows[3][1] = ""
	for _, tag := range t.TimeEntry.Tags {
		t.Rows[3][1] += tag.Name + " "
	}
}
