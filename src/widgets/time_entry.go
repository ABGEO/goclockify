package widgets

import (
	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
	"time"
)

type TimeEntryWidget struct {
	*w.Table
	TimeEntry TimeEntry
}

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

func (self *TimeEntryWidget) SetTimeEntry(timeEntry TimeEntry) {
	self.TimeEntry = timeEntry
	self.UpdateTable()
}

func (self *TimeEntryWidget) UpdateTable() {
	self.Rows[0][1] = self.TimeEntry.Description
	self.Rows[1][1] = self.TimeEntry.Project.Name
	self.Rows[2][1] = self.TimeEntry.Project.ClientName
	self.Rows[4][1] = self.TimeEntry.TimeInterval.Start.Format("01/02/2006 15:04:05")
	self.Rows[5][1] = self.TimeEntry.TimeInterval.End.Format("01/02/2006 15:04:05")

	if self.TimeEntry.TimeInterval.End.IsZero() {
		self.Rows[6][1] = "Running - " + time.Now().Sub(self.TimeEntry.TimeInterval.Start).String()
	} else {
		self.Rows[6][1] = self.TimeEntry.TimeInterval.End.Sub(self.TimeEntry.TimeInterval.Start).String()
	}

	self.Rows[3][1] = ""
	for _, tag := range self.TimeEntry.Tags {
		self.Rows[3][1] += tag.Name + " "
	}
}
