package widgets

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"strconv"
	"time"
)

type Tag struct {
	ID   string
	Name string
}

type Project struct {
	ID         string
	Name       string
	ClientName string
}

type TimeInterval struct {
	Start    time.Time
	End      time.Time
	Duration string
}

type TimeEntry struct {
	ID           string
	Description  string
	Tags         []Tag
	Billable     bool
	Project      Project
	TimeInterval TimeInterval
}

type TimeEntriesWidget struct {
	*Table
	TimeEntries []TimeEntry
}

func NewTimeEntriesWidget() *TimeEntriesWidget {
	self := &TimeEntriesWidget{
		Table: NewTable(),
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

func (self *TimeEntriesWidget) SetTimeEntries(timeEntries []TimeEntry) {
	self.TimeEntries = timeEntries
	self.entriesToRows()
	self.SelectedItem = ""
	self.ScrollTop()
}

func (self *TimeEntriesWidget) UpdateData(timeEntries []TimeEntry, workplace Workplace) {
	self.SetTimeEntries(timeEntries)
	self.Title = fmt.Sprintf(" %s - Time Entries ", workplace.Name)
}

func (self *TimeEntriesWidget) GetSelectedTimeEntry() (TimeEntry, error) {
	selectedIndex := self.Rows[self.SelectedRow][0]
	i, err := strconv.Atoi(selectedIndex)
	if err != nil {
		return TimeEntry{}, err
	}

	return self.TimeEntries[i], nil
}

func (self *TimeEntriesWidget) entriesToRows() {
	var timeEntries *[]TimeEntry
	timeEntries = &self.TimeEntries
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
	self.Rows = strings
}
