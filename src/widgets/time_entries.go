package widgets

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"time"
)

type TimeInterval struct {
	Start    time.Time
	End      time.Time
	Duration string
}

type TimeEntry struct {
	ID           string
	Description  string
	TagIds       []string
	Billable     bool
	ProjectId    string
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
	self.PadLeft = 2
	self.UniqueCol = 1
	self.Header = []string{"Description", "Duration", ""}
	self.ColResizer = func() {
		self.ColWidths = []int{ui.MaxInt(self.Inner.Dx()-26, 100), 15, 1}
	}

	return self
}

func (self *TimeEntriesWidget) SetTimeEntries(timeEntries []TimeEntry) {
	self.TimeEntries = timeEntries
	self.entriesToRows()
}

func (self *TimeEntriesWidget) UpdateData(timeEntries []TimeEntry, workplace Workplace) {
	self.SetTimeEntries(timeEntries)
	self.Title = fmt.Sprintf(" %s - Time Entries ", workplace.Name)
}

func (self *TimeEntriesWidget) entriesToRows() {
	var timeEntries *[]TimeEntry
	timeEntries = &self.TimeEntries
	strings := make([][]string, len(*timeEntries))
	for i, t := range *timeEntries {
		duration := t.TimeInterval.End.Sub(t.TimeInterval.Start)

		strings[i] = make([]string, 3)
		strings[i][0] = t.Description
		strings[i][1] = duration.String()

		if t.Billable {
			strings[i][2] = "$"
		} else {
			strings[i][2] = ""
		}
	}
	self.Rows = strings
}
